package collections

import (
	"strconv"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Create a new test runner we'll use to test all the
// modules in the collections package
func TestCollections(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Collections Suite")
}

var _ = Describe("Common Tests", func() {

	// Tests that calling ToDictionary with an empty map will result in a panic
	It("ToDictionary - Map is nil - Panic", func() {

		// Attempt to convert the list to a map; this should panic
		list := []int{1, 24, 3, 5}
		Expect(func() {
			ToDictionary(nil, list, func(i int) string {
				return strconv.FormatInt(int64(i), 10)
			}, true)
		}).Should(Panic())
	})

	// Tests that calling ToDictionary with a nil list will
	// result in no change to the map
	It("ToDictionary - List is nil - No work done", func() {

		// Attempt to convert a nil list to a map; nothing should happen
		mapping := make(map[string]int)
		ToDictionary(mapping, nil, func(i int) string {
			return strconv.FormatInt(int64(i), 10)
		}, true)

		// Verify that the map is still empty
		Expect(mapping).Should(BeEmpty())
	})

	// Tests that, if the ToDictionary function is called with a list that
	// contains elements resulting in collisions and that the overwrite
	// flag is true, then the map will contain the newer value
	It("ToDictionary - Overwrite true - Collisions overwritten", func() {

		// First, create our test list
		list := []int{1, 24, 3, 5, 16}

		// Next, convert the list to a map
		mapping := make(map[string]int)
		ToDictionary(mapping, list, func(i int) string {
			if i == 5 {
				return "3"
			}

			return strconv.FormatInt(int64(i), 10)
		}, true)

		// Finally, verify the values in the map
		Expect(mapping).Should(HaveLen(4))
		Expect(mapping["1"]).Should(Equal(1))
		Expect(mapping["3"]).Should(Equal(5))
		Expect(mapping["16"]).Should(Equal(16))
		Expect(mapping["24"]).Should(Equal(24))
	})

	// Tests that, if the ToDictionary function is called with a list that
	// contains elements resulting in collisions and that the overwrite
	// flag is false, then the map will contain the older value
	It("ToDictionary - Overwrite false - Collisions ignored", func() {

		// First, create our test list
		list := []int{1, 24, 3, 5, 16}

		// Next, convert the list to a map
		mapping := make(map[string]int)
		ToDictionary(mapping, list, func(i int) string {
			if i == 5 {
				return "3"
			}

			return strconv.FormatInt(int64(i), 10)
		}, false)

		// Finally, verify the values in the map
		Expect(mapping).Should(HaveLen(4))
		Expect(mapping["1"]).Should(Equal(1))
		Expect(mapping["3"]).Should(Equal(3))
		Expect(mapping["16"]).Should(Equal(16))
		Expect(mapping["24"]).Should(Equal(24))
	})

	// Tests that, if the AsSlice function is called with no data, then an
	// empty list will be returned
	It("AsSlice - No data provided - Empty list returned", func() {
		list := AsSlice[int]()
		Expect(list).Should(BeEmpty())
	})

	// Tests that, if the AsSlice function is called with data, then that data
	// will be added to a new slice of the same length that respects the ordering
	// of the data provided
	It("AsSlice - Data provided - Returned as list", func() {
		list := AsSlice(1, 2, 3, 10)
		Expect(list).Should(HaveLen(4))
		Expect(list).Should(Equal([]int{1, 2, 3, 10}))
	})

	// Tests that the Convert function will produce an empty list if called with no data
	It("Convert - No data provided - Empty list returned", func() {

		// Attempt to do the conversion with an empty list
		list := Convert(func(i int) string {
			return strconv.FormatInt(int64(i), 10)
		})

		// Verify that the list is empty
		Expect(list).Should(BeEmpty())
	})

	// Tests that the Convert function will produce a list of data where each item is the
	// result of an input to the convert function provided
	It("Convert - Data provided - Converted", func() {

		// Attempt to do the conversion with a non-empty list
		list := Convert(func(i int) string {
			return strconv.FormatInt(int64(i), 10)
		}, 1, 2, 3, 10)

		// Verify the converted data
		Expect(list).Should(HaveLen(4))
		Expect(list).Should(Equal([]string{"1", "2", "3", "10"}))
	})
})
