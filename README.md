<p align="center"><img src="https://raw.githubusercontent.com/go-composites/brand/main/social/golang-oop.png" alt="go-composites/typed" width="720"></p>

# go-composites/typed

A **type-safe, generics-based** variant of the go-composites building blocks.

The original go-composites repos (`result`, `array`, `optional`/`null`, …) carry
`interface{}` payloads. That buys *dynamic introspection* — a value can be
queried with `Kind()`, `RespondTo()`, reflective dispatch — but it pays for it
in **lost compile-time type safety**: `Result.Payload()` hands you an
`interface{}`, and a wrong type assertion compiles cleanly and only fails (panic
or `ok == false`) at run time.

This module recovers that safety with Go generics. There is **no `interface{}`
in the public API**. The compiler verifies that producers and consumers agree on
the payload type. It is a *parallel* module — it does not import or modify any of
the existing repos, and reimplements just enough to stand alone.

```
module github.com/go-composites/typed   (go 1.21+, CGO_ENABLED=0, zero deps)
```

## Packages

### `src/result` — `Result[T]`

Railway-oriented, statically-typed success/failure.

```go
type Result[T any] interface {
    HasError() bool
    Payload() (T, bool)
    Error() error
}

func Ok[T any](v T) Result[T]
func Err[T any](e error) Result[T]            // nil error normalised to a sentinel

func Map[T, U any](Result[T], func(T) U) Result[U]
func FlatMap[T, U any](Result[T], func(T) Result[U]) Result[U]
func AndThen[T, U any](Result[T], func(T) Result[U]) Result[U]   // alias for FlatMap
```

`Map` and `FlatMap` short-circuit on the failure branch — the error is
propagated and `f` is never called. `Ok` and `Err` are **real values, never
nil**, preserving the Null-Object spirit.

### `src/optional` — `Optional[T]`

The type-safe Null-Object. `None` is a concrete value, not a nil pointer.

```go
func Some[T any](v T) Optional[T]
func None[T any]() Optional[T]

func (Optional[T]) IsPresent() bool
func (Optional[T]) Get() (T, bool)
func (Optional[T]) OrElse(fallback T) T
func Map[T, U any](Optional[T], func(T) U) Optional[U]
```

### `src/slice` — `Slice[T]` (a typed Array)

A `[]T` with the usual combinators — no casts, no `interface{}`.

```go
func Of[T any](items ...T) Slice[T]
func Map[T, U any](Slice[T], func(T) U) Slice[U]
func Reduce[T, A any](Slice[T], initial A, combine func(A, T) A) A

func (Slice[T]) Filter(keep func(T) bool) Slice[T]
func (Slice[T]) Find(pred func(T) bool) optional.Optional[T]   // None on no match
func (Slice[T]) Any(pred func(T) bool) bool
func (Slice[T]) All(pred func(T) bool) bool
func (Slice[T]) Len() int
```

> Go methods cannot introduce new type parameters, so the operations whose
> result type differs from `T` (`Map`, `Reduce`) are free functions; the rest
> are methods.

## The trade-off (honest version)

| | dynamic go-composites (`interface{}`) | go-composites/typed (generics) |
|---|---|---|
| Type errors caught | at **run time** (panic / `ok=false`) | at **compile time** |
| `Payload()` | `interface{}` (needs a cast) | `(T, bool)` (no cast) |
| Reflective `Kind()` / `RespondTo()` | yes — values are introspectable | **no** — `T` is erased to a single monomorphised shape; there is nothing to interrogate |
| Heterogeneous containers | natural (`[]interface{}`) | needs a sum/`any` element type, losing the safety |
| Null-Object (`None`/`Err` never nil) | yes | yes (preserved) |

Generics are not strictly better — they are a **different point on the curve**.
You gain a compiler that rejects the mismatched-type bug the dynamic version can
only discover at run time. You lose the *reflective dynamism* that makes the
untyped composites composable across unknown types: there is no `Kind()` to ask
"what are you?", because the answer was fixed (and erased) at instantiation.

## Build / test

```sh
GOWORK=off CGO_ENABLED=0 go build ./... \
  && go vet ./... \
  && go test ./...

# 100% statement coverage on the library packages:
GOWORK=off CGO_ENABLED=0 go test -coverprofile=cover.out ./src/... \
  && go tool cover -func=cover.out | tail -1
```

## Demo

`go run .` exercises all three packages and contrasts them with the dynamic
variant (see `main.go`).

## License

BSD-3-Clause © 2026 the go-composites/typed authors.
