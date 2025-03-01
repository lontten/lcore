package lcutils

// BoolDiff
// list1-list2
// 布尔-差集,
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

// BoolUnion
// list1+list2
// 布尔-并,
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
