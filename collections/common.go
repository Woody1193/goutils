package collections

// ToDictionary converts a list of items to a map of items based on the output of
// a function that gets a key from each item and a Boolean value that determines
// whether conflicts should be overwritten
func ToDictionary[T any, U comparable](mapping map[U]T, list []T, keyer func(T) U, overwrite bool) {
	for _, item := range list {

		// Use the keyer to get a key from the item
		key := keyer(item)

		// If the mapping already contains the item then we'll either
		// ignore it, or if overwrite is true, we'll save the item to
		// the map; otherwise, save the item anyway
		if _, ok := mapping[key]; !ok || overwrite {
			mapping[key] = item
		}
	}
}

// AsSlice converts a parameterized list of items to a slice
func AsSlice[T any](data ...T) []T {
	return data
}

// Convert converts all the items in a list from a first type to
// a second type, using the function provided
func Convert[T any, U any](converter func(T) U, data ...T) []U {
	converted := make([]U, len(data))
	for i, item := range data {
		converted[i] = converter(item)
	}

	return converted
}
