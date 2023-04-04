package slices

import (
	"errors"
	"sort"
	"sync"
)

var RoutineCap int = 5

/*
Map every item of the slice using the mapFunc.
If the mapFunc returns a err for position i the new slice returns an empty value instead of the result at position i.
*/
func MapSlice[T, R any](slice []T, mapFunc func(T) (R, error)) ([]R, error) {
	if len(slice) == 0 {
		return []R{}, nil
	}
	var (
		wg                sync.WaitGroup
		countingSemaphore       = make(chan struct{}, RoutineCap)
		newBList                = make([]R, len(slice))
		errMutex                = make(chan struct{}, 1)
		err               error = nil
	)
	for i, elem := range slice {
		i := i
		elem := elem
		countingSemaphore <- struct{}{}
		wg.Add(1)
		go func(i int, elem T) {
			mappedElement, mapErr := mapFunc(elem)
			if mapErr != nil {
				errMutex <- struct{}{}
				if err == nil {
					err = mapErr
				} else {
					err = errors.Join(err, mapErr)
				}
				<-errMutex
			} else {
				newBList[i] = mappedElement
			}
			<-countingSemaphore
			wg.Done()
		}(i, elem)
	}
	wg.Wait()
	return newBList, err
}

// * We need concurrentItem to reorder the filtered slice
type concurrentItem[T any] struct {
	index int
	value T
}

/*
	Filter every item of the slice using the filterFunc.

If filterFunc returns an error the item is not included in the filtered slice.
*/
func FilterSlice[T any](slice []T, filterFunc func(T) (bool, error)) ([]T, error) {
	if len(slice) == 0 {
		return []T{}, nil
	}
	var (
		wg                  sync.WaitGroup
		countingSemaphore         = make(chan struct{}, RoutineCap)
		concurrentItemsChan       = make(chan concurrentItem[T], len(slice))
		errMutex                  = make(chan struct{}, 1)
		err                 error = nil
	)

	for i, elem := range slice {
		i := i
		elem := elem
		countingSemaphore <- struct{}{}
		wg.Add(1)
		go func(i int, elem T) {
			include, filterErr := filterFunc(elem)

			if filterErr != nil {
				errMutex <- struct{}{}
				if err == nil {
					err = filterErr
				} else {
					err = errors.Join(err, filterErr)
				}
				<-errMutex

			} else if include {
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
	return result, err
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

/*
Reduce slice to single value using the reduceFunc.
Use this if order of reduction matters.
If reduceFunc returns an error the result of reduceFunc is not included in the result.
*/
func OrderedReduceSlice[T, R any](slice []T, reduceFunc func(element T, basis R) (R, error), zero R) (R, error) {
	if len(slice) == 0 {
		return zero, nil
	}
	var (
		result       = zero
		err    error = nil
	)
	for _, element := range slice {
		newResult, reduceErr := reduceFunc(element, result)
		if reduceErr != nil {
			if err == nil {
				err = reduceErr
			} else {
				err = errors.Join(err, reduceErr)
			}
			continue
		}
		result = newResult
	}
	return result, err
}

/*
	Reduce slice to single value using the reduceFunc.

Use this order of reduction does not matter.
If reduceFunc returns an error the result of reduceFunc is not included in the result.
*/
func UnorderedReduceSlice[T, R any](slice []T, reduceFunc func(element T, basis R) (R, error), zero R) (R, error) {
	if len(slice) == 0 {
		return zero, nil
	}
	var (
		wg                sync.WaitGroup
		countingSemaphore       = make(chan struct{}, RoutineCap)
		mutex                   = make(chan struct{}, 1)
		result                  = zero
		errMutex                = make(chan struct{}, 1)
		err               error = nil
	)

	for _, elem := range slice {
		elem := elem
		countingSemaphore <- struct{}{}
		wg.Add(1)
		go func(element T) {
			mutex <- struct{}{}
			newResult, reduceErr := reduceFunc(element, result)
			if reduceErr != nil {
				errMutex <- struct{}{}
				if err == nil {
					err = reduceErr
				} else {
					err = errors.Join(err, reduceErr)
				}
				<-errMutex
			} else {
				result = newResult
			}
			<-mutex
			<-countingSemaphore
			wg.Done()
		}(elem)
	}
	wg.Wait()
	return result, err
}

/*
	Test if Every item fulfils isXFunc criteria.

If isXFunc returns an error the item is skipped for evaluation
*/
func EverySlice[T any](slice []T, isXFunc func(element T) (bool, error)) (bool, error) {
	if len(slice) == 0 {
		return false, nil
	}
	var (
		wg                sync.WaitGroup
		countingSemaphore       = make(chan struct{}, RoutineCap)
		result                  = true
		errMutex                = make(chan struct{}, 1)
		err               error = nil
	)

	for _, elem := range slice {
		elem := elem
		countingSemaphore <- struct{}{}
		wg.Add(1)
		go func(element T) {
			isX, isXErr := isXFunc(element)
			if isXErr != nil {
				errMutex <- struct{}{}
				if err == nil {
					err = isXErr
				} else {
					err = errors.Join(err, isXErr)
				}
				<-errMutex
			} else if result && !isX {
				result = false
			}
			<-countingSemaphore
			wg.Done()
		}(elem)
	}
	wg.Wait()
	return result, err
}

/*
Test if Some items fulfil isXFunc criteria.
If isXFunc returns an error the item is skipped for evaluation
*/
func SomeSlice[T any](slice []T, isXFunc func(element T) (bool, error)) (bool, error) {
	if len(slice) == 0 {
		return false, nil
	}
	var (
		wg       sync.WaitGroup
		control        = make(chan struct{}, RoutineCap)
		result         = false
		errMutex       = make(chan struct{}, 1)
		err      error = nil
	)
	for _, elem := range slice {
		elem := elem
		control <- struct{}{}
		wg.Add(1)
		go func(element T) {
			isX, isXErr := isXFunc(element)
			if isXErr != nil {
				errMutex <- struct{}{}
				if err == nil {
					err = isXErr
				} else {
					err = errors.Join(err, isXErr)
				}
				<-errMutex
			} else if !result && isX {
				result = true
			}
			<-control
			wg.Done()
		}(elem)
	}
	wg.Wait()
	return result, err
}

// Flat function :: [[]]a -> []a
func FlatSlice[a any](slices [][]a) (result []a) {
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return
}
