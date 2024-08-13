package pool

import (
	"errors"
	"sync"
)

var (
	bufferSize        int = 10
	defaultNumWorkers int = 10
)

// Run executes the task on a pool of workers of length numWorkers with the values
func Run[K any](numWorkers int, values []K, job func(index int, value K)) error {
	if numWorkers <= 0 {
		numWorkers = defaultNumWorkers
	}

	if len(values) == 0 {
		return errors.New("no values provided")
	}

	var wg sync.WaitGroup
	data := make(chan K, bufferSize)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, data, job, &wg)
	}

	for _, val := range values {
		data <- val
	}

	close(data)
	wg.Wait()

	return nil
}

// RunAndReturn executes the task on a pool of workers of length numWorkers with the values
// and returns the results.
// Keep in mind the results slice might not be in the same order
func RunAndReturn[K any, V any](numWorkers int, values []K, job func(index int, value K) V) ([]V, error) {
	if numWorkers <= 0 {
		numWorkers = defaultNumWorkers
	}

	if len(values) == 0 {
		return nil, errors.New("no values provided")
	}

	var wg sync.WaitGroup
	data := make(chan K, bufferSize)
	ret := make(chan V, bufferSize)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go workerWithReturn(i, data, ret, job, &wg)
	}

	for _, val := range values {
		data <- val
	}

	close(data)
	close(ret)
	wg.Wait()

	results := make([]V, len(values))
	for r := range ret {
		results = append(results, r)
	}

	return results, nil
}

// Sets the channels buffer size, default is 10
func SetBufferSize(size int) {
	if size < 0 {
		return
	}

	bufferSize = size
}

func worker[K any](index int, values <-chan K, job func(index int, value K), wg *sync.WaitGroup) {
	defer wg.Done()
	for val := range values {
		job(index, val)
	}
}

func workerWithReturn[K any, V any](index int, values <-chan K, results chan<- V, job func(index int, value K) V, wg *sync.WaitGroup) {
	defer wg.Done()
	for val := range values {
		results <- job(index, val)
	}
}
