package reqr

import (
	"reflect"
	"regexp"
	"testing"
)

type expectation struct {
	testing.TB

	got interface{}
}

func (e *expectation) Equals(v interface{}) {
	var (
		exp = v
		got = e.got
	)
	if !reflect.DeepEqual(exp, got) {
		e.Errorf("expected %s, got %s", exp, got)
	}
}

func (e *expectation) Contains(v interface{}) {
	switch t := e.got.(type) {
	case string:
		expr, ok := v.(string)
		if !ok {
			e.Fatal("%T: value type does not match", v)
		}
		reg, err := regexp.Compile(expr)
		if err != nil {
			e.Fatal(err)
		}
		if !reg.MatchString(t) {
			e.Errorf("expected %s to contain %s", t, expr)
		}

	default:
		e.Fatalf("%T: has no valid contains assertion for type")
	}
}
