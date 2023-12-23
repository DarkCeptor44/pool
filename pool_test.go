package pool

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testFolder = "temp"

func TestAction(t *testing.T) {
	arr := []string{}
	for i := 0; i < 5; i++ {
		arr = append(arr, fmt.Sprintf("file%d.txt", i+1))
	}

	p := NewPool(10, arr)
	err := p.Run(func(worker Worker, value string) error {
		t.Logf("Manipulating file %s by %s\n", value, worker.Name())
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

func TestReturn(t *testing.T) {
	arr := []string{}
	for i := 0; i < 5; i++ {
		arr = append(arr, fmt.Sprintf("file%d.txt", i+1))
	}

	p := NewReturnPool[string, bool](10, arr)
	res, err := p.Run(func(worker Worker, value string) (bool, error) {
		return strings.Contains(value, fmt.Sprint(worker)), nil
	})
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v\n", res)
}

func BenchmarkStd(b *testing.B) {
	b.Run("Read", wrapperForFiles(func(arr []string) {
		for _, file := range arr {
			f, err := os.Open(file)
			handleError(b, err)
			defer f.Close()

			var bb []byte
			_, err = f.Read(bb)
			handleError(b, err)
		}
	}))

	b.Run("Write", wrapperForFiles(func(arr []string) {
		for _, file := range arr {
			f, err := os.OpenFile(file, os.O_WRONLY, os.ModePerm)
			handleError(b, err)
			defer f.Close()

			_, err = f.WriteString("hi there")
			handleError(b, err)
		}
	}))
}

func BenchmarkAction(b *testing.B) {
	b.Run("Read", wrapperForFiles(func(arr []string) {
		p := NewPool(10, arr)
		err := p.Run(func(worker Worker, value string) error {
			f, err := os.Open(value)
			if err != nil {
				return err
			}
			defer f.Close()

			var bb []byte
			_, err = f.Read(bb)
			return err
		})
		handleError(b, err)
	}))

	b.Run("Write", wrapperForFiles(func(arr []string) {
		p := NewPool(10, arr)
		err := p.Run(func(worker Worker, value string) error {
			f, err := os.OpenFile(value, os.O_WRONLY, os.ModePerm)
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = f.WriteString("hi there")
			return err
		})
		handleError(b, err)
	}))
}

func BenchmarkReturn(b *testing.B) {
	b.Run("Read", wrapperForFiles(func(arr []string) {
		p := NewReturnPool[string, int](10, arr)
		res, err := p.Run(func(worker Worker, value string) (int, error) {
			f, err := os.Open(value)
			if err != nil {
				return 0, nil
			}
			defer f.Close()

			var bb []byte
			n, err := f.Read(bb)
			return n, err
		})
		handleError(b, err)
		if len(res) == 0 {
			b.Log("no results")
			b.FailNow()
		}
	}))
}

func wrapperForFiles(a func(arr []string)) func(b *testing.B) {
	return func(b *testing.B) {
		arr := make([]string, 0)
		os.Mkdir(testFolder, os.ModePerm)

		for i := 0; i < 10; i++ {
			name := fmt.Sprintf("file%d.txt", i+1)
			path := filepath.Join(testFolder, name)
			f, err := os.Create(path)
			handleError(b, err)
			err = f.Truncate(1024 * 1024)
			handleError(b, err)
			f.Close()
			arr = append(arr, path)
		}

		for i := 0; i < b.N; i++ {
			a(arr)
		}

		for _, file := range arr {
			err := os.Remove(file)
			handleError(b, err)
		}
	}
}

func handleError(b *testing.B, err error) {
	if err != nil {
		b.Log(err)
		b.FailNow()
	}
}
