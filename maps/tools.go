package maps

func CloneAndClear[T ~map[K]V, K comparable, V any](m T) map[K]V {
	var res = make(map[K]V, len(m))
	for k, v := range m {
		delete(m, k)
		res[k] = v
	}
	return res
}
