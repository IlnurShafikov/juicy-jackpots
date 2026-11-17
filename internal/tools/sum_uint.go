package tools

// SumUint64 calculates the sum of all elements in a uint64 slice.
// The function iterates over all elements of the slice and returns their total sum.
//
// Parameters:
//   - slice: slice of unsigned 64-bit integers to sum
//
// Returns:
//   - sum of all elements in the slice
//   - 0 if the slice is empty
func SumUint64(slice []uint64) uint64 {
	var sum uint64

	for _, v := range slice {
		sum += v
	}

	return sum
}
