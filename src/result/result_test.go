package result

import (
	"errors"
	"strconv"
	"testing"
)

func TestOkPayloadAndError(t *testing.T) {
	r := Ok(42)
	if r.HasError() {
		t.Fatal("Ok should not have an error")
	}
	v, ok := r.Payload()
	if !ok || v != 42 {
		t.Fatalf("Payload = (%d,%v), want (42,true)", v, ok)
	}
	if r.Error() != nil {
		t.Fatalf("Ok.Error = %v, want nil", r.Error())
	}
}

func TestErrPayloadAndError(t *testing.T) {
	e := errors.New("boom")
	r := Err[int](e)
	if !r.HasError() {
		t.Fatal("Err should have an error")
	}
	v, ok := r.Payload()
	if ok || v != 0 {
		t.Fatalf("Payload = (%d,%v), want (0,false)", v, ok)
	}
	if r.Error() != e {
		t.Fatalf("Error = %v, want %v", r.Error(), e)
	}
}

func TestErrNilIsNormalised(t *testing.T) {
	r := Err[string](nil)
	if !r.HasError() {
		t.Fatal("Err(nil) must still be a failure")
	}
	if r.Error() == nil {
		t.Fatal("Err(nil) must carry a non-nil sentinel error")
	}
}

func TestMapSuccessAndFailure(t *testing.T) {
	ok := Map(Ok(3), func(n int) string { return strconv.Itoa(n * 2) })
	if v, _ := ok.Payload(); v != "6" {
		t.Fatalf("Map(Ok) = %q, want %q", v, "6")
	}

	e := errors.New("x")
	fail := Map(Err[int](e), func(n int) string { return strconv.Itoa(n) })
	if !fail.HasError() || fail.Error() != e {
		t.Fatalf("Map(Err) should propagate the error %v, got %v", e, fail.Error())
	}
	if _, ok := fail.Payload(); ok {
		t.Fatal("Map(Err) must not carry a payload")
	}
}

func parsePositive(s string) Result[int] {
	n, err := strconv.Atoi(s)
	if err != nil {
		return Err[int](err)
	}
	if n <= 0 {
		return Err[int](errors.New("not positive"))
	}
	return Ok(n)
}

func TestFlatMapAndThen(t *testing.T) {
	good := FlatMap(Ok("21"), parsePositive)
	if v, _ := good.Payload(); v != 21 {
		t.Fatalf("FlatMap good = %d, want 21", v)
	}

	// failure in the inner function
	bad := AndThen(Ok("-1"), parsePositive)
	if !bad.HasError() {
		t.Fatal("AndThen should surface the inner failure")
	}

	// short-circuit: upstream failure never calls f
	e := errors.New("upstream")
	called := false
	short := FlatMap(Err[string](e), func(s string) Result[int] {
		called = true
		return Ok(0)
	})
	if called {
		t.Fatal("FlatMap must not call f on the failure branch")
	}
	if short.Error() != e {
		t.Fatalf("short.Error = %v, want %v", short.Error(), e)
	}
}
