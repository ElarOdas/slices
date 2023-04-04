package slices_test

import (
	"errors"
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
		result, _ := slices.MapSlice(intSource, func(i int) (string, error) { return strconv.Itoa(i), nil })

		if !reflect.DeepEqual(result, strConvMapTarget) {
			t.Fail()
			break
		}
	}
}

func TestConcurrencyFilter(t *testing.T) {
	for i := 0; i < 100; i++ {
		result, _ := slices.FilterSlice(intSource, func(i int) (bool, error) {
			return i >= 30, nil
		})
		if !reflect.DeepEqual(result, filter30Target) {
			t.Fail()
			break
		}
	}
}
func TestConcurrencyReduce(t *testing.T) {
	for i := 0; i < 100; i++ {
		result, _ := slices.UnorderedReduceSlice(intSource, func(i int, base int) (int, error) {
			return base + i, nil
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
		result, _ := slices.EverySlice(intSource, func(i int) (bool, error) {
			return i > 0, nil
		})
		if !result {
			t.Fail()
			break
		}
	}
	for i := 0; i < 50; i++ {
		// 50 Tests for false
		result, _ := slices.EverySlice(intSource, func(i int) (bool, error) {
			return i > 50, nil
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
		result, _ := slices.SomeSlice(intSource, func(i int) (bool, error) {
			return i > 50, nil
		})
		if !result {
			t.Fail()
			break
		}
	}
	for i := 0; i < 50; i++ {
		// 50 Tests for true
		result, _ := slices.SomeSlice(intSource, func(i int) (bool, error) {
			return i < 0, nil
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
	result, _ := slices.MapSlice(complexSource, func(elem testStruct) (string, error) {
		return fmt.Sprintf("%s%d", elem.name, (elem.a + elem.b)), nil
	})
	if !reflect.DeepEqual(result, complexMapTarget) {
		t.Fail()
	}
}
func TestComplexFilter(t *testing.T) {
	result, _ := slices.FilterSlice(complexSource, func(elem testStruct) (bool, error) {
		return elem.a+elem.b >= 80, nil
	})
	if !reflect.DeepEqual(result, complexFilterTarget) {
		t.Fail()
	}
}

func TestComplexOrderedReduce(t *testing.T) {
	result, _ := slices.OrderedReduceSlice(complexSource, func(elem testStruct, base int) (int, error) {
		return base + elem.a + elem.b, nil
	}, 0)
	if result != complexReduceTarget {
		t.Fail()
	}
}
func TestComplexUnorderedReduce(t *testing.T) {
	result, _ := slices.UnorderedReduceSlice(complexSource, func(elem testStruct, base int) (int, error) {
		return base + elem.a + elem.b, nil
	}, 0)
	if result != complexReduceTarget {
		t.Fail()
	}
}
func TestComplexEvery(t *testing.T) {
	resultTrue, _ := slices.EverySlice(complexSource, func(elem testStruct) (bool, error) {
		return elem.a > 0, nil
	})

	resultFalse, _ := slices.EverySlice(complexSource, func(elem testStruct) (bool, error) {
		return elem.a > 50, nil
	})
	if resultFalse || !resultTrue {
		t.Fail()
	}

}
func TestComplexSome(t *testing.T) {
	resultTrue, _ := slices.SomeSlice(complexSource, func(elem testStruct) (bool, error) {
		return elem.a > 50, nil
	})
	resultFalse, _ := slices.SomeSlice(complexSource, func(elem testStruct) (bool, error) {
		return elem.a < 0, nil
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
	result, _ := slices.MapSlice(emptyIntSlice, func(i int) (string, error) {
		return "", nil
	})
	if !reflect.DeepEqual(result, emptyStringSlice) {
		t.Fail()
	}
}
func TestEmptyFilter(t *testing.T) {
	result, _ := slices.FilterSlice(emptyIntSlice, func(i int) (bool, error) { return true, nil })
	if !reflect.DeepEqual(result, emptyIntSlice) {
		t.Fail()
	}
}
func TestEmptyOrderedReduce(t *testing.T) {
	result, _ := slices.OrderedReduceSlice(emptyIntSlice, func(i int, basis int) (int, error) {
		return basis + i, nil
	}, 0)
	if result != 0 {
		t.Fail()
	}
}
func TestEmptyUnorderedReduce(t *testing.T) {
	result, _ := slices.UnorderedReduceSlice(emptyIntSlice, func(i int, basis int) (int, error) {
		return basis + i, nil
	}, 0)
	if result != 0 {
		t.Fail()
	}
}
func TestEmptyEvery(t *testing.T) {
	result, _ := slices.EverySlice(emptyIntSlice, func(i int) (bool, error) {
		return i > 50, nil
	})
	if result {
		t.Fail()
	}
}
func TestEmptySome(t *testing.T) {
	result, _ := slices.SomeSlice(emptyIntSlice, func(i int) (bool, error) {
		return i > 50, nil
	})
	if result {
		t.Fail()
	}
}

// ? Tests for errors should be included

var errorSource = []string{
	"a",
	"b",
	"x",
	"56",
	"2",
}

var errorMapTarget = func() []int {
	var result = make([]int, len(errorSource))
	for i, element := range errorSource {
		num, err := strconv.Atoi(element)
		if err != nil {
			continue
		}
		result[i] = num
	}
	return result
}()

var errorFilterTarget = func() []string {
	var result []string
	for _, element := range errorSource {
		num, err := strconv.Atoi(element)
		if err != nil || num <= 3 {
			continue
		}

		result = append(result, element)
	}
	return result
}()

var errorReduceTarget = func() int {
	var result int
	for _, element := range errorSource {
		num, err := strconv.Atoi(element)
		if err != nil {
			continue
		}
		result += num
	}
	return result
}()

func TestMapError(t *testing.T) {
	result, err := slices.MapSlice(errorSource, strconv.Atoi)

	var numErr *strconv.NumError

	if !errors.As(err, &numErr) {
		t.Fail()
	}
	if !reflect.DeepEqual(result, errorMapTarget) {
		t.Fail()
	}
}

func TestFilterError(t *testing.T) {
	result, err := slices.FilterSlice(errorSource, func(s string) (bool, error) {
		i, err := strconv.Atoi(s)
		if err != nil {
			return false, err
		}
		return i > 3, err
	})
	var numErr *strconv.NumError
	if !errors.As(err, &numErr) {
		t.Fail()
	}
	if !reflect.DeepEqual(result, errorFilterTarget) {
		t.Fail()
	}
}

func TestOrderedReduceError(t *testing.T) {
	result, err := slices.OrderedReduceSlice(errorSource, func(s string, base int) (int, error) {
		i, err := strconv.Atoi(s)
		if err != nil {
			return base, err
		}
		return base + i, err
	}, 0)
	var numErr *strconv.NumError
	if !errors.As(err, &numErr) {
		t.Fail()
	}
	if !reflect.DeepEqual(result, errorReduceTarget) {
		t.Error(result, errorReduceTarget)
	}
}
func TestUnorderedReduceError(t *testing.T) {
	result, err := slices.UnorderedReduceSlice(errorSource, func(s string, base int) (int, error) {
		i, err := strconv.Atoi(s)
		if err != nil {
			return base, err
		}
		return base + i, err
	}, 0)
	var numErr *strconv.NumError
	if !errors.As(err, &numErr) {
		t.Fail()
	}
	if !reflect.DeepEqual(result, errorReduceTarget) {
		t.Fail()
	}
}

func TestEveryError(t *testing.T) {
	var numErr *strconv.NumError
	// Test true
	resultTrue, err := slices.EverySlice(errorSource, func(s string) (bool, error) {
		i, err := strconv.Atoi(s)
		if err != nil {
			return false, err
		}
		return i > 1, err
	})
	if !errors.As(err, &numErr) {
		t.Fail()
	}
	if !resultTrue {
		t.Fail()
	}
	// Test false
	resultFalse, err := slices.EverySlice(errorSource, func(s string) (bool, error) {
		i, err := strconv.Atoi(s)
		if err != nil {
			return false, err
		}
		return i > 3, err
	})
	if !errors.As(err, &numErr) {
		t.Fail()
	}
	if resultFalse {
		t.Fail()
	}

}
func TestSomeError(t *testing.T) {
	var numErr *strconv.NumError
	// Test true
	resultTrue, err := slices.SomeSlice(errorSource, func(s string) (bool, error) {
		i, err := strconv.Atoi(s)
		if err != nil {
			return false, err
		}
		return i > 1, err
	})
	if !errors.As(err, &numErr) {
		t.Fail()
	}
	if !resultTrue {
		t.Fail()
	}
	// Test false
	resultFalse, err := slices.SomeSlice(errorSource, func(s string) (bool, error) {
		i, err := strconv.Atoi(s)
		if err != nil {
			return false, err
		}
		return i < -1, err
	})
	if !errors.As(err, &numErr) {
		t.Fail()
	}
	if resultFalse {
		t.Fail()
	}
}
