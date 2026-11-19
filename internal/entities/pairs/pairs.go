// Package pairs provides a simple ordered key-value list of integers
// that is read from configuration files.
package pairs

// Pair is a single key-value entry.
type Pair struct {
	// Key is the entry identifier (for example a symbol ID).
	Key int
	// Value is the value associated with the key.
	Value int
}

// Pairs is an ordered list of key-value entries.
type Pairs []Pair

// GetValueByKey returns the value of the first entry with the given key.
//
// Parameters:
//   - key: key to search for
//
// Returns:
//   - value of the entry
//   - true if the key was found, false otherwise
func (kv Pairs) GetValueByKey(key int) (int, bool) {
	for i := range kv {
		curr := kv[i]

		if curr.Key == key {
			return curr.Value, true
		}
	}

	return 0, false
}

// GetIDByValue returns the key of the first entry with the given value.
//
// Parameters:
//   - value: value to search for
//
// Returns:
//   - key of the entry
//   - true if the value was found, false otherwise
func (kv Pairs) GetIDByValue(value int) (int, bool) {
	for i := range kv {
		curr := kv[i]

		if curr.Value == value {
			return curr.Key, true
		}
	}

	return 0, false
}

// GetByIndex returns the entry at the given position.
//
// Parameters:
//   - index: zero-based position in the list
//
// Returns:
//   - entry at the position
//   - false if the index is out of range
func (kv Pairs) GetByIndex(index int) (Pair, bool) {
	if index < 0 || index >= len(kv) {
		return Pair{Key: 0, Value: 0}, false
	}

	return kv[index], true
}
