package slices

func Filter[T any](s []T, keep func(T) bool) []T {
	var result []T

	for _, v := range s {
		if keep(v) {
			result = append(result, v)
		}
	}

	return result
}
