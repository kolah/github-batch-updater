package slices

func Map[T any, U any](s []T, f func(T) U) []U {
	result := make([]U, 0, len(s))

	for _, v := range s {
		result = append(result, f(v))
	}

	return result
}
