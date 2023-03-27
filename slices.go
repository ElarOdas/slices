package slices

import (
	"sort"
	"sync"
)

var routineCap int = 5

func SetRoutineCap(newCap int) {
	routineCap = newCap
}

// Map
func MapSlice[T, R any](Tslice []T, f func(T) R) []R {
	if len(Tslice) == 0 {
		return []R{}
	}
	control := make(chan struct{}, routineCap)
	var wg sync.WaitGroup
	newBList := make([]R, len(Tslice))

	for i, elem := range Tslice {
		control <- struct{}{}
		wg.Add(1)
		go func(i int, elem T) {
			newBList[i] = f(elem)
			<-control
			wg.Done()
		}(i, elem)
	}
	wg.Wait()
	return newBList
}

type concurrentSliceItem[T any] struct {
	index int
	value T
}

// Filter
func FilterSlice[T any](Tslice []T, f func(T) bool) []T {
	if len(Tslice) == 0 {
		return []T{}
	}
	var wg sync.WaitGroup
	control := make(chan struct{}, routineCap)
	out := make(chan concurrentSliceItem[T], len(Tslice))

	for i, elem := range Tslice {
		control <- struct{}{}
		wg.Add(1)
		go func(i int, elem T) {
			if f(elem) {
				out <- concurrentSliceItem[T]{i, elem}
			}
			<-control
			wg.Done()
		}(i, elem)
	}
	wg.Wait()
	close(out)

	var outList = chanToList(out)

	sort.Slice(outList, func(i, j int) bool {
		return outList[i].index < outList[j].index
	})
	result := make([]T, len(outList))
	for i := range outList {
		result[i] = outList[i].value
	}
	return result
}

func chanToList[T any](c chan T) []T {
	var result = make([]T, 0)
	for elem := range c {
		result = append(result, elem)
	}
	return result
}

// Reduce
func OrderdReduceSlice[T, R any](Tslice []T, reduction func(element T, basis R) R, zero R) R {
	if len(Tslice) == 0 {
		return zero
	}
	result := zero
	for _, element := range Tslice {
		result = reduction(element, result)
	}
	return result
}

func UnorderedReduceSlice[T, R any](Tslice []T, reduction func(element T, basis R) R, zero R) R {
	if len(Tslice) == 0 {
		return zero
	}
	result := zero
	var wg sync.WaitGroup
	control := make(chan struct{}, routineCap)
	mutex := make(chan struct{}, 1)

	for _, element := range Tslice {
		control <- struct{}{}
		wg.Add(1)
		go func(element T) {
			mutex <- struct{}{}
			result = reduction(element, result)
			<-mutex
			<-control
			wg.Done()
		}(element)
	}
	wg.Wait()
	return result
}

// ? Consider coupling every and some
// Every
func EverySlice[T any](Tslice []T, tester func(element T) bool) bool {
	result := true
	var wg sync.WaitGroup
	control := make(chan struct{}, routineCap)
	for _, element := range Tslice {
		control <- struct{}{}
		wg.Add(1)
		go func(element T) {
			if result && !tester(element) {
				result = false
			}
			<-control
			wg.Done()
		}(element)
	}
	wg.Wait()
	return result
}

// Some
func SomeSlice[T any](Tslice []T, tester func(element T) bool) bool {
	result := false
	var wg sync.WaitGroup
	control := make(chan struct{}, routineCap)
	for _, element := range Tslice {
		control <- struct{}{}
		wg.Add(1)
		go func(element T) {
			if !result && tester(element) {
				result = true
			}
			<-control
			wg.Done()
		}(element)
	}
	wg.Wait()
	return result
}

// Flat function :: [[]]a -> []a
func FlatSlice[a any](slices [][]a) (result []a) {
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return
}
