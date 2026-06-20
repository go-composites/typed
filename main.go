// Command typed is a demo contrasting the generics-based composites in this
// module with the dynamic (interface{}-based) go-composites repos.
//
// The headline difference is invisible at run time: it is what the compiler
// rejects. With the dynamic Result, Payload() returns interface{} and a wrong
// type assertion compiles fine and only panics (or returns ok=false) at run
// time. Here, Result[int].Payload() returns (int, bool) — feeding it where a
// string is expected does not compile at all.
package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-composites/typed/src/optional"
	"github.com/go-composites/typed/src/result"
	"github.com/go-composites/typed/src/slice"
)

// parseQty is a fallible operation returning a railway-typed Result[int].
func parseQty(s string) result.Result[int] {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return result.Err[int](fmt.Errorf("parse %q: %w", s, err))
	}
	if n < 0 {
		return result.Err[int](errors.New("quantity must be non-negative"))
	}
	return result.Ok(n)
}

func main() {
	fmt.Println("== go-composites/typed — compile-time-safe composites ==")

	// --- Result: railway-oriented, statically typed ----------------------
	// Chain parse -> double -> render, type-checked end to end.
	good := result.Map(
		result.AndThen(result.Ok("21"), parseQty),
		func(n int) string { return fmt.Sprintf("doubled=%d", n*2) },
	)
	if v, ok := good.Payload(); ok {
		fmt.Println("Result(ok):  ", v) // Result(ok):   doubled=42
	}

	bad := result.AndThen(result.Ok("oops"), parseQty)
	fmt.Println("Result(err): ", bad.HasError(), "->", bad.Error())

	// The dynamic version would let this compile:
	//     payload := dynResult.Payload()        // interface{}
	//     s := payload.(string)                 // panics if it was an int
	// Here, good.Payload() is (string, bool) — the mistake is a compile error.

	// --- Optional: the type-safe Null-Object -----------------------------
	first := slice.Of(3, 5, 8, 13).Find(func(n int) bool { return n > 7 })
	fmt.Println("Find>7:      ", first.OrElse(-1)) // 8

	missing := slice.Of(1, 2).Find(func(n int) bool { return n > 99 })
	fmt.Println("Find>99:     ", missing.IsPresent(), "OrElse=", missing.OrElse(-1))

	upper := optional.Map(optional.Some("hi"), strings.ToUpper)
	fmt.Println("Optional.Map:", upper.OrElse("?")) // HI

	// --- Slice: a typed Array, no interface{} casts ----------------------
	nums := slice.Of(1, 2, 3, 4, 5, 6)
	evens := nums.Filter(func(n int) bool { return n%2 == 0 })
	labels := slice.Map(evens, func(n int) string { return fmt.Sprintf("#%d", n) })
	total := slice.Reduce(nums, 0, func(acc, n int) int { return acc + n })

	fmt.Println("evens:       ", []int(evens))
	fmt.Println("labels:      ", []string(labels))
	fmt.Println("sum:         ", total)
	fmt.Println("any even:    ", nums.Any(func(n int) bool { return n%2 == 0 }))
	fmt.Println("all < 10:    ", nums.All(func(n int) bool { return n < 10 }))
}
