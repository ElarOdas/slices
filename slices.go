package slices

import (
	"sort"
	"sync"
)

var RoutineCap int = 5

// Map every item of the slice using the mapFunc
func MapSlice[T, R any](slice []T, mapFunc func(T) R) []R {
	if len(slice) == 0 {
		return []R{}
	}
	var (
		wg                sync.WaitGroup
		countingSemaphore = make(chan struct{}, RoutineCap)
		newBList          = make([]R, len(slice))
	)
	for i, elem := range slice {
		i := i
		elem := elem
		countingSemaphore <- struct{}{}
		wg.Add(1)
		go func(i int, elem T) {
			newBList[i] = mapFunc(elem)
			<-countingSemaphore
			wg.Done()
		}(i, elem)
	}
	wg.Wait()
	return newBList
}

type concurrentItem[T any] struct {
	index int
	value T
}

// * We need concurrentItem to reorder the filtered slice
// Filter every item of the slice using the filterFunc
func FilterSlice[T any](slice []T, filterFunc func(T) bool) []T {
	if len(slice) == 0 {
		return []T{}
	}
	var (
		wg                  sync.WaitGroup
		countingSemaphore   = make(chan struct{}, RoutineCap)
		concurrentItemsChan = make(chan concurrentItem[T], len(slice))
	)

	for i, elem := range slice {
		i := i
		elem := elem
		countingSemaphore <- struct{}{}
		wg.Add(1)
		go func(i int, elem T) {

			if filterFunc(elem) {
				concurrentItemsChan <- concurrentItem[T]{i, elem}
			}
			<-countingSemaphore
			wg.Done()
		}(i, elem)
	}
	wg.Wait()
	close(concurrentItemsChan)

	var concurrentSlice = chanToConcurrentSlice(concurrentItemsChan)
	result := concurrentSliceToSlice(concurrentSlice)
	return result
}

// Extract the slice of concurrentItems from chan
func chanToConcurrentSlice[T any](c chan T) []T {
	var result = make([]T, 0)
	for elem := range c {
		result = append(result, elem)
	}
	return result
}

// Reorder the slice and extract
func concurrentSliceToSlice[T any](concurrentSlice []concurrentItem[T]) []T {
	result := make([]T, len(concurrentSlice))

	sort.Slice(concurrentSlice, func(i, j int) bool {
		return concurrentSlice[i].index < concurrentSlice[j].index
	})
	for i := range concurrentSlice {
		result[i] = concurrentSlice[i].value
	}
	return result
}

// Reduce slice to single value using the reduceFunc
// Order of reduction matters
func OrderedReduceSlice[T, R any](slice []T, reduceFunc func(element T, basis R) R, zero R) R {
	if len(slice) == 0 {
		return zero
	}
	result := zero
	for _, element := range slice {
		result = reduceFunc(element, result)
	}
	return result
}

// Reduce slice to single value using the reduceFunc
// Order of reduction does not matter
func UnorderedReduceSlice[T, R any](slice []T, reduceFunc func(element T, basis R) R, zero R) R {
	if len(slice) == 0 {
		return zero
	}
	var (
		wg                sync.WaitGroup
		countingSemaphore = make(chan struct{}, RoutineCap)
		mutex             = make(chan struct{}, 1)
		result            = zero
	)

	for _, elem := range slice {
		elem := elem
		countingSemaphore <- struct{}{}
		wg.Add(1)
		go func(element T) {
			mutex <- struct{}{}
			result = reduceFunc(element, result)
			<-mutex
			<-countingSemaphore
			wg.Done()
		}(elem)
	}
	wg.Wait()
	return result
}

// Every item fulfils isXFunc criteria
func EverySlice[T any](slice []T, isXFunc func(element T) bool) bool {
	if len(slice) == 0 {
		return false
	}
	var (
		wg                sync.WaitGroup
		countingSemaphore = make(chan struct{}, RoutineCap)
		result            = true
	)

	for _, elem := range slice {
		elem := elem
		countingSemaphore <- struct{}{}
		wg.Add(1)
		go func(element T) {
			if result && !isXFunc(element) {
				result = false
			}
			<-countingSemaphore
			wg.Done()
		}(elem)
	}
	wg.Wait()
	return result
}

// Some items fulfil isXFunc criteria
func SomeSlice[T any](slice []T, isXFunc func(element T) bool) bool {
	if len(slice) == 0 {
		return false
	}
	var (
		wg      sync.WaitGroup
		control = make(chan struct{}, RoutineCap)
		result  = false
	)
	for _, elem := range slice {
		elem := elem
		control <- struct{}{}
		wg.Add(1)
		go func(element T) {
			if !result && isXFunc(element) {
				result = true
			}
			<-control
			wg.Done()
		}(elem)
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
