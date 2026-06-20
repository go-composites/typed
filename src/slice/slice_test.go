package slice

import (
	"strconv"
	"testing"
)

func eq[T comparable](t *testing.T, got, want []T) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d (got %v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("index %d = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestOfAndLen(t *testing.T) {
	s := Of(1, 2, 3)
	if s.Len() != 3 {
		t.Fatalf("Len = %d, want 3", s.Len())
	}
}

func TestMap(t *testing.T) {
	out := Map(Of(1, 2, 3), func(n int) string { return strconv.Itoa(n * 10) })
	eq(t, out, []string{"10", "20", "30"})
}

func TestFilter(t *testing.T) {
	evens := Of(1, 2, 3, 4, 5).Filter(func(n int) bool { return n%2 == 0 })
	eq(t, evens, []int{2, 4})

	none := Of(1, 3).Filter(func(n int) bool { return n%2 == 0 })
	if none.Len() != 0 {
		t.Fatalf("expected empty, got %v", none)
	}
}

func TestReduce(t *testing.T) {
	sum := Reduce(Of(1, 2, 3, 4), 0, func(acc, n int) int { return acc + n })
	if sum != 10 {
		t.Fatalf("Reduce sum = %d, want 10", sum)
	}
	// accumulator type differs from element type
	joined := Reduce(Of(1, 2, 3), "", func(acc string, n int) string {
		return acc + strconv.Itoa(n)
	})
	if joined != "123" {
		t.Fatalf("Reduce join = %q, want %q", joined, "123")
	}
}

func TestFind(t *testing.T) {
	found := Of(1, 2, 3).Find(func(n int) bool { return n > 1 })
	if v, ok := found.Get(); !ok || v != 2 {
		t.Fatalf("Find = (%d,%v), want (2,true)", v, ok)
	}

	missing := Of(1, 2, 3).Find(func(n int) bool { return n > 9 })
	if missing.IsPresent() {
		t.Fatal("Find on no match must return None")
	}
}

func TestAny(t *testing.T) {
	if !Of(1, 2, 3).Any(func(n int) bool { return n == 2 }) {
		t.Fatal("Any should be true")
	}
	if Of(1, 2, 3).Any(func(n int) bool { return n == 9 }) {
		t.Fatal("Any should be false")
	}
}

func TestAll(t *testing.T) {
	if !Of(2, 4, 6).All(func(n int) bool { return n%2 == 0 }) {
		t.Fatal("All should be true")
	}
	if Of(2, 3, 6).All(func(n int) bool { return n%2 == 0 }) {
		t.Fatal("All should be false")
	}
	if !Of[int]().All(func(n int) bool { return false }) {
		t.Fatal("All on empty should be vacuously true")
	}
}
