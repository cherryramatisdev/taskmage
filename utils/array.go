package utils

// Flatten takes a slice of slices of any type and returns a flattened slice.
func Flatten[T any](slices [][]T) []T {
	var result []T
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]struct{})
	var result []T
	for _, item := range slice {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
