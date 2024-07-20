package pool_test

import (
	"math/rand"
	"slices"
	"testing"
	"time"

	"github.com/DarkCeptor44/pool"
)

func TestNormal(t *testing.T) {
	wrapTest(t, "Normal", func(i []int, f func(int)) {
		for _, v := range i {
			f(v)
		}
	})
}

func TestPool(t *testing.T) {
	t.Run("NoValues", func(t *testing.T) {
		err := pool.Run(10, []int{}, func(i int, v int) {
			t.Log(i, v)
		})
		if err == nil {
			t.Fatal("expected an error")
		}
	})

	wrapTest(t, "Pool", func(i []int, f func(int)) {
		pool.Run(10, i, func(_, v int) {
			f(v)
		})
	})
}

func TestPoolWithReturn(t *testing.T) {
	t.Run("NoValues", func(t *testing.T) {
		_, err := pool.RunAndReturn(10, []int{}, func(i int, v int) int {
			return v * v
		})
		if err == nil {
			t.Fatal("expected an error")
		}
	})
}

func BenchmarkNormal(b *testing.B) {
	wrap(b, func(i []int, f func(int)) {
		for _, v := range i {
			f(v)
		}
	})
}

func BenchmarkPool(b *testing.B) {
	wrap(b, func(i []int, f func(int)) {
		pool.Run(10, i, func(_, v int) {
			f(v)
		})
	})
}

func wrap(b *testing.B, f func([]int, func(int))) {
	vals := []int{3, 2, 1, 5, 6, 7, 8, 9, 10, 32}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		f(vals, func(v int) {
			if !slices.Contains(vals, v) {
				b.Fatal("unexpected value")
			}
		})
	}
}

func wrapTest(t *testing.T, name string, f func([]int, func(int))) {
	num := 5

	t.Run(name, func(t *testing.T) {
		var sum time.Duration

		for i := 0; i < num; i++ {
			vals := randSlice(10000, 100000)
			start := time.Now()
			f(vals, func(v int) {
				if !slices.Contains(vals, v) {
					t.Fatal("unexpected value")
				}
			})
			sum += time.Since(start)
		}

		average := sum / time.Duration(num)
		t.Logf("'%s' Took %s\n", name, average)
	})
}

func randSlice(m, n int) []int {
	vals := make([]int, n)
	for i := 0; i < n; i++ {
		vals[i] = rand.Intn(m)
	}
	return vals
}
