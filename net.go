package lcore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func Get[T any](targetURL string) (T, error) {
	statusCode, result, err := GetBase[T](targetURL)
	if err != nil {
		return result, err
	}
	if statusCode != http.StatusOK {
		return result, fmt.Errorf("请求失败，状态码: %d", statusCode)
	}
	return result, nil
}
func GetBase[T any](targetURL string) (int, T, error) {
	var result T
	resp, err := http.Get(targetURL)
	if err != nil {
		return 0, result, fmt.Errorf("发送 GET 请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		// 读取错误响应体（可选，用于更详细的错误信息）
		errBody, _ := io.ReadAll(resp.Body)
		return 0, result, fmt.Errorf("请求失败，状态码: %d，响应: %s", resp.StatusCode, string(errBody))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, result, fmt.Errorf("读取响应体失败: %w", err)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, result, fmt.Errorf("响应体 JSON 反序列化失败: %w", err)
	}
	return 0, result, nil
}
func PostJson[T any](targetURL string, data any) (T, error) {
	statusCode, result, err := PostJsonBase[T](targetURL, data)
	if err != nil {
		return result, err
	}
	if statusCode != http.StatusOK {
		return result, fmt.Errorf("请求失败，状态码: %d", statusCode)
	}
	return result, nil
}
func PostJsonBase[T any](targetURL string, data any) (int, T, error) {
	var result T
	var jsonBody []byte
	var err error

	if data == nil {
		jsonBody = []byte("{}")
	} else {
		jsonBody, err = json.Marshal(data)
		if err != nil {
			return 0, result, fmt.Errorf("请求体 JSON 序列化失败: %w", err)
		}
	}

	resp, err := http.Post(targetURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return 0, result, fmt.Errorf("发送 PostJson 请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		// 读取错误响应体（可选，用于更详细的错误信息）
		errBody, _ := io.ReadAll(resp.Body)
		return 0, result, fmt.Errorf("请求失败，状态码: %d，响应: %s", resp.StatusCode, string(errBody))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, result, fmt.Errorf("读取响应体失败: %w", err)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, result, fmt.Errorf("响应体 JSON 反序列化失败: %w", err)
	}
	return resp.StatusCode, result, nil
}

func PostForm[T any](targetURL string, data url.Values) (T, error) {
	statusCode, result, err := PostFormBase[T](targetURL, data)
	if err != nil {
		return result, err
	}
	if statusCode != http.StatusOK {
		return result, fmt.Errorf("请求失败，状态码: %d", statusCode)
	}
	return result, nil
}
func PostFormBase[T any](targetURL string, data url.Values) (int, T, error) {
	var result T

	if data == nil {
		data = url.Values{}
	}

	resp, err := http.PostForm(targetURL, data)
	if err != nil {
		return 0, result, fmt.Errorf("发送 PostForm 请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		// 读取错误响应体（可选，用于更详细的错误信息）
		errBody, _ := io.ReadAll(resp.Body)
		return 0, result, fmt.Errorf("请求失败，状态码: %d，响应: %s", resp.StatusCode, string(errBody))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, result, fmt.Errorf("读取响应体失败: %w", err)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, result, fmt.Errorf("响应体 JSON 反序列化失败: %w", err)
	}
	return resp.StatusCode, result, nil
}
