package lcutils

// BoolDiff
// list1-list2
// 布尔-差集,
// 有去重逻辑
// 移除list1 中 所有的 list2 元素
func BoolDiff[T comparable](list1, list2 []T) []T {
	// 创建集合存储list1的元素
	m := make(map[T]struct{})
	for _, e := range list1 {
		m[e] = struct{}{}
	}

	// 过滤list2中不在集合里的元素
	for _, e := range list2 {
		if _, ok := m[e]; ok {
			delete(m, e)
		}
	}
	list1 = make([]T, 0, len(m))
	for e := range m {
		list1 = append(list1, e)
	}
	return list1
}

// BoolEq
// 布尔-相等
// 有去重逻辑
func BoolEq[T comparable](list1, list2 []T) bool {
	return len(BoolDiff(list1, list2)) == 0 && len(BoolDiff(list2, list1)) == 0
}

// BoolUnion
// list1+list2
// 布尔-并,
// 有去重逻辑
// list1 和 list2 元素，并集
func BoolUnion[T comparable](list1, list2 []T) []T {
	// 创建集合存储list1的元素
	m := make(map[T]struct{})
	for _, e := range list1 {
		m[e] = struct{}{}
	}
	for _, e := range list2 {
		m[e] = struct{}{}
	}

	list1 = make([]T, 0, len(m))
	for e := range m {
		list1 = append(list1, e)
	}
	return list1
}

// BoolIntersection
// list1，list2相同的元素
// 布尔-交集,
// 有去重逻辑
// list1 和 list2 元素，交集
func BoolIntersection[T comparable](list1, list2 []T) []T {
	// 创建集合存储list1的元素
	m := make(map[T]struct{})
	for _, e := range list1 {
		m[e] = struct{}{}
	}
	list1 = make([]T, 0, len(m))
	for _, e := range list2 {
		if _, ok := m[e]; ok {
			list1 = append(list1, e)
		}
	}
	return list1
}

// 去重函数，适用于任何可比较的类型
func RemoveDuplicates[T comparable](slice []T) []T {
	// 创建一个 map 用于存储已出现的元素
	seen := make(map[T]bool)
	result := []T{}

	// 遍历切片，将未重复的元素添加到结果中
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

// 集合中是否包含item
func ListHas[T comparable](slice []T, item T) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

type ListToolBuilder struct {
	list   []any
	hasAll []any
	hasAny []any
	notAll []any
}

func ListTool(list ...any) *ListToolBuilder {
	return &ListToolBuilder{
		list:   list,
		hasAll: []any{},
		hasAny: []any{},
		notAll: []any{},
	}
}

func (t *ListToolBuilder) HasAll(list ...any) *ListToolBuilder {
	t.hasAll = append(t.hasAll, list...)
	return t
}
func (t *ListToolBuilder) HasAny(list ...any) *ListToolBuilder {
	t.hasAny = append(t.hasAny, list...)
	return t
}
func (t *ListToolBuilder) NotAll(list ...any) *ListToolBuilder {
	t.notAll = append(t.notAll, list...)
	return t
}
func (t *ListToolBuilder) Check() bool {
	if len(t.hasAll) > 0 {
		var c = len(BoolIntersection(t.list, t.hasAll)) == len(t.hasAll)
		if !c {
			return false
		}
	}
	if len(t.hasAny) > 0 {
		var c = len(BoolIntersection(t.list, t.hasAny)) > 0
		if !c {
			return false
		}
	}
	if len(t.notAll) > 0 {
		var c = len(BoolIntersection(t.list, t.notAll)) == 0
		if !c {
			return false
		}

	}
	return true
}
