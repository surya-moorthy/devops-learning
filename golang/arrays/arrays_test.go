package arrays

import (
	"reflect"
	"testing"
)

func TestSum(t *testing.T) {
	t.Run("collections of 5 strings", func(t *testing.T) {
		numbers := []int{1, 5, 6, 2, 7}
		got := Sum(numbers)
		want := 21

		if got != want {
			t.Errorf("got %d want %d given %v", got, want, numbers)
		}
	})
	t.Run("collection of multiple", func(t *testing.T) {
		numbers := []int{1, 4, 25, 6, 2}
		got := Sum(numbers)
		want := 38
		if got != want {
			t.Errorf("got %d want %d given %v", got, want, numbers)
		}
	})
}

func TestSumAll(t *testing.T) {
	got := SumAll([]int{1, 5}, []int{6, 3})
	want := []int{7, 9}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %d want %d ", got, want)
	}
}

func BenchmarkSum(b *testing.B) {
	for b.Loop() {
		Sum([]int{14, 3, 6, 2, 6})
	}
}
