package utils

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
