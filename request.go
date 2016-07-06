package reqr

import (
	"io"
	"net/http"
	"reflect"
	"testing"
)

type request struct {
	testing.TB

	req *http.Request
}

var _ Request = &request{}

func (r *request) Header(k string) Expectation {
	v, ok := r.req.Header[k]
	if !ok {
		r.Errorf("expected Request Header: %s", k)
	}

	return &expectation{
		TB: r.TB,

		got: v,
	}
}

func (r *request) Body(tr ...TransformFunc) Expectation {
	var (
		got interface{}

		lentr = len(tr)
	)
	if lentr == 0 {
		body := r.req.Body
		defer body.Close()

		var b []byte
		n, err := io.ReadFull(body, b)
		if err != nil {
			r.Fatal(err)
		}
		got = string(b[:n])
	} else {
		var (
			in interface{} = r.req.Body

			i = 0
			j = lentr
		)
		for ; i < j; i++ {
			vo_v := reflect.ValueOf(tr[i])

			// TODO first transform fn should always be func(io.ReadCloser)

			ret := vo_v.Call([]reflect.Value{reflect.ValueOf(in)})

			// TODO heck for errors

			got = ret[0].Interface()
		}

	}

	return &expectation{
		TB: r.TB,

		got: got,
	}
}

func (r *request) Host() Expectation {
	return &expectation{
		TB: r.TB,

		got: r.req.Host,
	}
}

func (r *request) Request() *http.Request {
	return r.req
}
