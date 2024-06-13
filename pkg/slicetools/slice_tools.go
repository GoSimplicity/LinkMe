package slicetools

func Map[T any, U any](slice []T, mapper func(int, T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = mapper(i, v)
	}
	return result
}
