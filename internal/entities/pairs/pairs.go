package pairs

// Pair holds a single integer key-value association parsed from a CSV row.
type Pair struct {
	Key   int
	Value int
}

// Pairs is an ordered slice of Pair entries. Linear search is used for all
// lookups, which is acceptable given the small sizes typical in game configs.
type Pairs []Pair

// GetValueByKey performs a linear search and returns the Value of the first
// Pair whose Key equals key.
//
// Returns:
//   - the matching value and true if found
//   - 0 and false if no Pair has the given key
func (kv Pairs) GetValueByKey(key int) (int, bool) {
	for i := range kv {
		curr := kv[i]

		if curr.Key == key {
			return curr.Value, true
		}
	}

	return 0, false
}

// GetIDByValue performs a linear search and returns the Key of the first
// Pair whose Value equals value. Useful for reverse-lookup of jackpot type IDs.
//
// Returns:
//   - the matching key and true if found
//   - 0 and false if no Pair has the given value
func (kv Pairs) GetIDByValue(value int) (int, bool) {
	for i := range kv {
		curr := kv[i]

		if curr.Value == value {
			return curr.Key, true
		}
	}

	return 0, false
}

// GetByIndex returns the Pair at the given zero-based index.
//
// Returns:
//   - the Pair and true if index is within bounds
//   - a zero-value Pair and false if index is negative or >= len(kv)
func (kv Pairs) GetByIndex(index int) (Pair, bool) {
	if index < 0 || index >= len(kv) {
		return Pair{Key: 0, Value: 0}, false
	}

	return kv[index], true
}
