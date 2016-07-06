package reqr

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type response struct {
	testing.TB

	req  *http.Request
	resp *httptest.ResponseRecorder
}

var _ Response = &response{}

func (r *response) Status() Expectation {
	return &expectation{
		TB: r.TB,

		got: r.resp.Code,
	}
}

func (r *response) Header(k string) Expectation {
	return &expectation{
		TB: r.TB,

		got: r.resp.Header().Get(k),
	}
}

// Body can be transformed a number of times using a TransformFunc signature
//
// The first transformation signature argument must be `io.Reader` as that is
// what the httptest.ResponseRecorder's Body field implements
func (r *response) Body(tr ...TransformFunc) Expectation {
	var (
		got interface{}

		lentr = len(tr)
	)
	if lentr == 0 {
		got = r.resp.Body.String()
	} else {
		var (
			in interface{} = r.resp.Body

			i = 0
			j = lentr
		)
		for ; i < j; i++ {
			tr_func := reflect.ValueOf(tr[i])

			// TODO first transform fn should always be func(io.Reader)

			ret := tr_func.Call([]reflect.Value{reflect.ValueOf(in)})

			in = ret[0].Interface()
			if err := ret[1].Interface(); err != nil {
				r.Errorf("transform error: %s", err)

				break
			}
		}

		got = in
	}

	return &expectation{
		TB: r.TB,

		got: got,
	}
}

func (r *response) Request() Request {
	return &request{
		TB: r.TB,

		req: r.req,
	}
}
