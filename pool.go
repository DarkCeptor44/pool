package pool

import (
	"fmt"
	"sync"
)

// Pool represents a pool instance
type Pool[K any] struct {
	NumWorkers int
	workers    []Worker
	data       []K
}

type ReturnPool[K any, V any] struct {
	NumWorkers int
	workers    []Worker
	data       []K
}

type Worker int

func (w Worker) Name() string {
	return fmt.Sprintf("worker%d", w+1)
}

// ActionTask represents a function that is executed by each worker but doesnt return anything.
// worker represents the number of the worker in the pool, starting from 1
type ActionTask[K any] func(worker Worker, value K) error

type ReturnTask[K any, V any] func(worker Worker, value K) (V, error)

// NewPool creates a new Pool where K is the type given to the workers.
// If number of workers is 0, 10 will be used
func NewPool[K any](numWorkers int, data []K) *Pool[K] {
	return &Pool[K]{
		NumWorkers: IfThenElse(numWorkers <= 0, 10, numWorkers),
		workers:    makeWorkers(numWorkers),
		data:       data,
	}
}

// NewReturnPool creates a new Pool where K is the type given to the workers and V is the type returned from the workers.
// If number of workers is 0, 10 will be used
func NewReturnPool[K any, V any](numWorkers int, data []K) *ReturnPool[K, V] {
	return &ReturnPool[K, V]{
		NumWorkers: IfThenElse(numWorkers <= 0, 10, numWorkers),
		workers:    makeWorkers(numWorkers),
		data:       data,
	}
}

// Executes the task on the pool
func (p *Pool[K]) Run(task ActionTask[K]) error {
	var wg sync.WaitGroup
	values := make(chan K)
	errors := make(chan error)
	num := min(p.NumWorkers, len(p.data))

	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(workerIndex int) {
			defer wg.Done()
			for value := range values {
				worker := p.workers[workerIndex]
				err := task(worker, value)
				if err != nil {
					errors <- err
				}
			}
		}(i)
	}

	for _, value := range p.data {
		values <- value
	}

	go func() {
		close(values)
		wg.Wait()
		close(errors)
	}()

	for err := range errors {
		if err != nil {
			return err
		}
	}

	return nil
}

// Executes the task on the pool
func (p *ReturnPool[K, V]) Run(task ReturnTask[K, V]) ([]V, error) {
	var wg sync.WaitGroup
	values := make(chan K)
	errors := make(chan error)
	results := make(chan V)
	num := min(p.NumWorkers, len(p.data))

	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(workerIndex int) {
			defer wg.Done()
			for value := range values {
				worker := p.workers[workerIndex]
				r, err := task(worker, value)
				if err != nil {
					errors <- err
				}
				results <- r
			}
		}(i)
	}

	for _, value := range p.data {
		values <- value
	}

	go func() {
		close(values)
		wg.Wait()
		close(errors)
		close(results)
	}()

	res := make([]V, 0)
	for r := range results {
		res = append(res, r)
	}

	for err := range errors {
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func makeWorkers(n int) []Worker {
	workers := make([]Worker, 0)
	for i := 0; i < n; i++ {
		workers = append(workers, Worker(i))
	}
	return workers
}

// Returns a if condition is true, b otherwise.
// Adapted from https://github.com/shomali11/util/blob/master/xconditions/xconditions.go#L12.
// Wanted something that uses generics
func IfThenElse[K any](condition bool, a, b K) K {
	if condition {
		return a
	}
	return b
}
