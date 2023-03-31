package slices_test

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/ElarOdas/slices"
)

// Test for Concurrency aka Run 100x: Test if there are unseen Race conditions

var (
	intSource        = []int{44, 24, 57, 19, 64, 90, 9, 37, 39, 47, 33, 54, 61, 45, 11, 17, 21, 8, 62, 23, 1, 87, 93, 60, 91, 81, 13, 18, 97, 72, 55, 92, 47, 45, 4, 30, 92, 22, 16, 98, 34, 80, 56, 95, 60, 69, 57, 71, 25, 32, 16, 60, 61, 37, 62, 83, 90, 57, 60, 54, 84, 49, 64, 50, 22, 81, 4, 50, 17, 44, 67, 89, 23, 27, 71, 22, 29, 80, 52, 73, 36, 30, 29, 15, 67, 78, 8, 51, 85, 44, 6, 68, 63, 73, 74, 89, 41, 2, 29, 41, 60, 73}
	strConvMapTarget = func() []string {
		result := []string{}
		for _, elem := range intSource {
			result = append(result, strconv.Itoa(elem))
		}
		return result
	}()
	filter30Target = func() []int {
		result := []int{}
		for _, elem := range intSource {
			if elem >= 30 {
				result = append(result, elem)
			}
		}
		return result
	}()
	reduceTarget = func() int {
		result := 0
		for _, elem := range intSource {
			result = result + elem

		}
		return result
	}()
)

func TestConcurrencyMap(t *testing.T) {
	for i := 0; i < 100; i++ {
		result := slices.MapSlice(intSource, strconv.Itoa)
		if !reflect.DeepEqual(result, strConvMapTarget) {
			t.Fail()
			break
		}
	}
}

func TestConcurrencyFilter(t *testing.T) {
	for i := 0; i < 100; i++ {
		result := slices.FilterSlice(intSource, func(i int) bool {
			return (i >= 30)
		})
		if !reflect.DeepEqual(result, filter30Target) {
			t.Fail()
			break
		}
	}
}
func TestConcurrencyReduce(t *testing.T) {
	for i := 0; i < 100; i++ {
		result := slices.UnorderedReduceSlice(intSource, func(i int, base int) int {
			return base + i
		}, 0)
		if result != reduceTarget {
			t.Fail()
			break
		}
	}
}

func TestConcurrencyEvery(t *testing.T) {
	for i := 0; i < 50; i++ {
		// 50 Tests for true
		result := slices.EverySlice(intSource, func(i int) bool {
			return i > 0
		})
		if !result {
			t.Fail()
			break
		}
	}
	for i := 0; i < 50; i++ {
		// 50 Tests for false
		result := slices.EverySlice(intSource, func(i int) bool {
			return i > 50
		})
		if result {
			t.Fail()
			break
		}
	}
}
func TestConcurrencySome(t *testing.T) {
	for i := 0; i < 50; i++ {
		// 50 Tests for true
		result := slices.SomeSlice(intSource, func(i int) bool {
			return i > 50
		})
		if !result {
			t.Fail()
			break
		}
	}
	for i := 0; i < 50; i++ {
		// 50 Tests for true
		result := slices.SomeSlice(intSource, func(i int) bool {
			return i < 0
		})
		if result {
			t.Fail()
			break
		}
	}
}

// Test for T = struct: Test what happens if T is a complex structure

type testStruct struct {
	name string
	a    int
	b    int
}

var (
	complexSource = []testStruct{
		{name: "Test0", a: 77, b: 56},
		{name: "Test1", a: 3, b: 80},
		{name: "Test2", a: 38, b: 39},
		{name: "Test3", a: 46, b: 52},
		{name: "Test4", a: 12, b: 47},
		{name: "Test5", a: 77, b: 28},
		{name: "Test6", a: 54, b: 0},
		{name: "Test7", a: 94, b: 53},
		{name: "Test8", a: 67, b: 90},
		{name: "Test9", a: 37, b: 56},
	}

	complexFlattenSource = [][]testStruct{
		{{name: "Test0", a: 77, b: 56}},
		{{name: "Test1", a: 3, b: 80}},
		{{name: "Test2", a: 38, b: 39}},
		{{name: "Test3", a: 46, b: 52}},
		{{name: "Test4", a: 12, b: 47}},
		{{name: "Test5", a: 77, b: 28}},
		{{name: "Test6", a: 54, b: 0}},
		{{name: "Test7", a: 94, b: 53}},
		{{name: "Test8", a: 67, b: 90}},
		{{name: "Test9", a: 37, b: 56}},
	}

	complexMapTarget = func() []string {
		result := []string{}
		for _, elem := range complexSource {
			result = append(result, fmt.Sprintf("%s%d", elem.name, (elem.a+elem.b)))
		}
		return result
	}()
	complexFilterTarget = func() []testStruct {
		result := []testStruct{}
		for _, elem := range complexSource {
			if elem.a+elem.b >= 80 {
				result = append(result, elem)
			}
		}
		return result
	}()

	complexReduceTarget = func() int {
		result := 0
		for _, elem := range complexSource {
			result = result + elem.a + elem.b
		}
		return result
	}()
)

func TestComplexMap(t *testing.T) {
	result := slices.MapSlice(complexSource, func(elem testStruct) string {
		return fmt.Sprintf("%s%d", elem.name, (elem.a + elem.b))
	})
	if !reflect.DeepEqual(result, complexMapTarget) {
		t.Fail()
	}
}
func TestComplexFilter(t *testing.T) {
	result := slices.FilterSlice(complexSource, func(elem testStruct) bool {
		return elem.a+elem.b >= 80
	})
	if !reflect.DeepEqual(result, complexFilterTarget) {
		t.Fail()
	}
}

func TestComplexOrderedReduce(t *testing.T) {
	result := slices.OrderedReduceSlice(complexSource, func(elem testStruct, base int) int {
		return base + elem.a + elem.b
	}, 0)
	if result != complexReduceTarget {
		t.Fail()
	}
}
func TestComplexUnorderedReduce(t *testing.T) {
	result := slices.UnorderedReduceSlice(complexSource, func(elem testStruct, base int) int {
		return base + elem.a + elem.b
	}, 0)
	if result != complexReduceTarget {
		t.Fail()
	}
}
func TestComplexEvery(t *testing.T) {
	resultTrue := slices.EverySlice(complexSource, func(elem testStruct) bool {
		return elem.a > 0
	})

	resultFalse := slices.EverySlice(complexSource, func(elem testStruct) bool {
		return elem.a > 50
	})
	if resultFalse || !resultTrue {
		t.Fail()
	}

}
func TestComplexSome(t *testing.T) {
	resultTrue := slices.SomeSlice(complexSource, func(elem testStruct) bool {
		return elem.a > 50
	})
	resultFalse := slices.SomeSlice(complexSource, func(elem testStruct) bool {
		return elem.a < 0
	})

	if resultFalse || !resultTrue {
		t.Fail()
	}
}

func TestFlatten(t *testing.T) {
	result := slices.FlatSlice(complexFlattenSource)

	if !reflect.DeepEqual(result, complexSource) {
		t.Fail()
	}
}

// Test for empty slice

var (
	emptyIntSlice    = []int{}
	emptyStringSlice = []string{}
)

func TestEmptyMap(t *testing.T) {
	result := slices.MapSlice(emptyIntSlice, func(i int) string {
		return ""
	})
	if !reflect.DeepEqual(result, emptyStringSlice) {
		t.Fail()
	}
}
func TestEmptyFilter(t *testing.T) {
	result := slices.FilterSlice(emptyIntSlice, func(i int) bool { return true })
	if !reflect.DeepEqual(result, emptyIntSlice) {
		t.Fail()
	}
}
func TestEmptyOrderedReduce(t *testing.T) {
	result := slices.OrderedReduceSlice(emptyIntSlice, func(i int, basis int) int {
		return basis + i
	}, 0)
	if result != 0 {
		t.Fail()
	}
}
func TestEmptyUnorderedReduce(t *testing.T) {
	result := slices.UnorderedReduceSlice(emptyIntSlice, func(i int, basis int) int {
		return basis + i
	}, 0)
	if result != 0 {
		t.Fail()
	}
}
func TestEmptyEvery(t *testing.T) {
	result := slices.SomeSlice(emptyIntSlice, func(i int) bool {
		return i > 50
	})
	if result {
		t.Fail()
	}
}
func TestEmptySome(t *testing.T) {
	result := slices.SomeSlice(emptyIntSlice, func(i int) bool {
		return i > 50
	})
	if result {
		t.Fail()
	}
}

// ? Tests for errors should be included
