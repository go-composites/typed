package optional

import (
	"strconv"
	"testing"
)

func TestSomePresent(t *testing.T) {
	o := Some(7)
	if !o.IsPresent() {
		t.Fatal("Some should be present")
	}
	v, ok := o.Get()
	if !ok || v != 7 {
		t.Fatalf("Get = (%d,%v), want (7,true)", v, ok)
	}
	if o.OrElse(99) != 7 {
		t.Fatalf("OrElse on Some = %d, want 7", o.OrElse(99))
	}
}

func TestNoneAbsent(t *testing.T) {
	o := None[int]()
	if o.IsPresent() {
		t.Fatal("None should be absent")
	}
	v, ok := o.Get()
	if ok || v != 0 {
		t.Fatalf("Get = (%d,%v), want (0,false)", v, ok)
	}
	if o.OrElse(99) != 99 {
		t.Fatalf("OrElse on None = %d, want 99", o.OrElse(99))
	}
}

func TestMap(t *testing.T) {
	some := Map(Some(5), func(n int) string { return strconv.Itoa(n) })
	if v, ok := some.Get(); !ok || v != "5" {
		t.Fatalf("Map(Some) = (%q,%v), want (\"5\",true)", v, ok)
	}

	none := Map(None[int](), func(n int) string { return strconv.Itoa(n) })
	if none.IsPresent() {
		t.Fatal("Map(None) must stay absent")
	}
}
