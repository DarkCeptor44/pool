package pool

import (
	"errors"
	"sync"
)

var (
	BufferSize        int = 10
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
	data := make(chan K, BufferSize)

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

func worker[K any](index int, values <-chan K, job func(index int, value K), wg *sync.WaitGroup) {
	for val := range values {
		job(index, val)
		wg.Done()
	}
}
