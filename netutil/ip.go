package netutil

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
)

// IPConfig 配置结构
type IPConfig struct {
	// 信任的代理服务器CIDR列表
	TrustedProxies []string
	// 信任的代理IP列表（已解析的net.IPNet）
	trustedNetworks []*net.IPNet
	// 是否启用严格模式（只返回验证通过的IP）
	StrictMode bool
	// 支持的IP头部优先级
	IPHeaders []string
	// 私有网络检查（新增）
	AllowPrivate bool
	once         sync.Once
	initErr      error
}

// DefaultIPConfig 默认配置
var DefaultIPConfig = &IPConfig{
	TrustedProxies: []string{
		"10.0.0.0/8",     // 私有网络
		"172.16.0.0/12",  // 私有网络
		"192.168.0.0/16", // 私有网络
		"127.0.0.0/8",    // 环回地址
		"fc00::/7",       // 私有IPv6网络
		"::1/128",        // IPv6环回
	},
	StrictMode:   false,
	AllowPrivate: true,
	IPHeaders: []string{
		"X-Real-IP",
		"X-Forwarded-For",
		"CF-Connecting-IP",    // Cloudflare
		"True-Client-IP",      // Akamai和Cloudflare
		"X-Cluster-Client-IP", // Rackspace LB, Riverbed Stingray
	},
}

// 自定义错误类型
var (
	ErrInvalidIP      = errors.New("invalid IP address")
	ErrUntrustedProxy = errors.New("untrusted proxy")
)

// init 初始化配置，解析CIDR
func (c *IPConfig) init() error {
	c.once.Do(func() {
		c.trustedNetworks = make([]*net.IPNet, 0, len(c.TrustedProxies))
		for _, proxy := range c.TrustedProxies {
			if _, network, err := net.ParseCIDR(proxy); err == nil {
				c.trustedNetworks = append(c.trustedNetworks, network)
			} else {
				c.initErr = fmt.Errorf("invalid CIDR %q: %w", proxy, err)
				return
			}
		}
	})
	return c.initErr
}

// isTrustedProxy 检查IP是否在信任的代理列表中
func (c *IPConfig) isTrustedProxy(ipStr string) bool {
	if err := c.init(); err != nil {
		return false
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	for _, network := range c.trustedNetworks {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

// isPrivateIP 检查是否为私有IP（新增）
func (c *IPConfig) isPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	// 简化私有IP检查，使用配置的信任网络
	return c.isTrustedProxy(ipStr)
}

// RealIP 生产环境获取真实IP
func RealIP(r *http.Request, config *IPConfig) (string, error) {
	if config == nil {
		config = DefaultIPConfig
	}

	if err := config.init(); err != nil {
		return "", err
	}

	var candidateIPs []string
	remoteIP := extractIPFromRemoteAddr(r.RemoteAddr)

	// 验证远程IP格式
	if remoteIP != "" && net.ParseIP(remoteIP) == nil {
		if config.StrictMode {
			return "", fmt.Errorf("%w: %s", ErrInvalidIP, remoteIP)
		}
		remoteIP = ""
	}

	// 收集所有候选IP
	for _, header := range config.IPHeaders {
		if value := r.Header.Get(header); value != "" {
			ips := parseIPHeader(value)
			// 验证每个IP的格式
			validIPs := make([]string, 0, len(ips))
			for _, ip := range ips {
				if net.ParseIP(ip) != nil {
					validIPs = append(validIPs, ip)
				}
			}
			candidateIPs = append(candidateIPs, validIPs...)
		}
	}

	// 选择最可信的IP
	ip := selectTrustedIP(candidateIPs, remoteIP, config)

	// 严格模式下验证最终IP
	if config.StrictMode {
		if ip == "" || net.ParseIP(ip) == nil {
			return "", ErrInvalidIP
		}
		if !config.AllowPrivate && config.isPrivateIP(ip) {
			return "", fmt.Errorf("private IP not allowed: %s", ip)
		}
	}

	return ip, nil
}

// parseIPHeader 解析IP头部，处理逗号分隔的多个IP
func parseIPHeader(headerValue string) []string {
	headerValue = strings.TrimSpace(headerValue)
	if headerValue == "" {
		return nil
	}

	var ips []string
	parts := strings.Split(headerValue, ",")

	for _, part := range parts {
		ip := strings.TrimSpace(part)
		if ip != "" {
			ips = append(ips, ip)
		}
	}

	return ips
}

// extractIPFromRemoteAddr 从RemoteAddr中提取IP
func extractIPFromRemoteAddr(remoteAddr string) string {
	if remoteAddr == "" {
		return ""
	}

	// 处理带端口的情况
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		// 可能没有端口，直接尝试解析
		if net.ParseIP(remoteAddr) != nil {
			return remoteAddr
		}
		return ""
	}

	return ip
}

// selectTrustedIP 选择最可信的IP
func selectTrustedIP(candidateIPs []string, remoteIP string, config *IPConfig) string {
	// 如果没有候选IP，直接返回remoteIP
	if len(candidateIPs) == 0 {
		return remoteIP
	}

	// 如果远程IP是信任的代理，返回第一个候选IP
	if config.isTrustedProxy(remoteIP) {
		return candidateIPs[0]
	}

	// 否则从候选IP中寻找第一个非信任代理的IP
	for i := len(candidateIPs) - 1; i >= 0; i-- {
		ip := candidateIPs[i]
		if !config.isTrustedProxy(ip) {
			return ip
		}
	}

	// 如果所有候选IP都是信任的代理，返回最后一个
	return candidateIPs[len(candidateIPs)-1]
}

// RealIPSimple 简化版本，使用默认配置
func RealIPSimple(r *http.Request) string {
	ip, _ := RealIP(r, DefaultIPConfig)
	return ip
}

// 中间件版本
func IPMiddleware(config *IPConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, err := RealIP(r, config)
			if err != nil && config != nil && config.StrictMode {
				http.Error(w, "Invalid client IP", http.StatusForbidden)
				return
			}

			// 将IP添加到请求上下文
			ctx := r.Context()
			// 使用自定义类型避免键冲突
			type ipKey struct{}
			ctx = context.WithValue(ctx, ipKey{}, ip)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// GetIPFromContext 从上下文中获取IP
func GetIPFromContext(ctx context.Context) string {
	type ipKey struct{}
	if ip, ok := ctx.Value(ipKey{}).(string); ok {
		return ip
	}
	return ""
}

// 工具函数：验证IP格式
func IsValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// 工具函数：获取IP版本（4或6）
func GetIPVersion(ip string) int {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return 0
	}
	if parsedIP.To4() != nil {
		return 4
	}
	return 6
}

// 新增：检查IP是否在CIDR范围内
func IsIPInCIDR(ipStr, cidrStr string) (bool, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false, ErrInvalidIP
	}

	_, network, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return false, err
	}

	return network.Contains(ip), nil
}
