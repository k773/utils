package maps

func CloneAndClear[T ~map[K]V, K comparable, V any](m T) map[K]V {
	var res = make(map[K]V, len(m))
	for k, v := range m {
		delete(m, k)
		res[k] = v
	}
	return res
}

func ValuesAndClear[T ~map[K]V, K comparable, V any](m T) []V {
	var res = make([]V, len(m))
	var i int
	for k, v := range m {
		delete(m, k)
		res[i] = v
		i++
	}
	return res
}

func KeysAndClear[T ~map[K]V, K comparable, V any](m T) []K {
	var res = make([]K, len(m))
	var i int
	for k := range m {
		delete(m, k)
		res[i] = k
		i++
	}
	return res
}

func Clear[T ~map[K]V, K comparable, V any](m T) {
	for k := range m {
		delete(m, k)
	}
}

func Clone[T ~map[K]V, K comparable, V any](m T) map[K]V {
	var res = make(map[K]V, len(m))
	for k, v := range m {
		res[k] = v
	}
	return res
}

func Values[T ~map[K]V, K comparable, V any](m T) []V {
	var res = make([]V, len(m))
	var i int
	for _, v := range m {
		res[i] = v
		i++
	}
	return res
}

func Keys[T ~map[K]V, K comparable, V any](m T) []K {
	var res = make([]K, len(m))
	var i int
	for k := range m {
		res[i] = k
		i++
	}
	return res
}
