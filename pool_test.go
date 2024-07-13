package pool_test

import (
	"slices"
	"testing"

	"github.com/DarkCeptor44/pool"
)

func TestPool(t *testing.T) {
	t.Run("NoValues", func(t *testing.T) {
		err := pool.Run(10, []int{}, func(i int, v int) {
			t.Log(i, v)
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
		pool.Run(5, i, func(_, v int) {
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
