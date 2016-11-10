package reqr

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type ReqFunc func(*http.Request)

type Reqr interface {
	GET(string, ...ReqFunc) Response
	POST(string, interface{}, ...ReqFunc) Response
	PUT(string, interface{}, ...ReqFunc) Response
	DELETE(string, interface{}, ...ReqFunc) Response
	OPTIONS(string, interface{}, ...ReqFunc) Response
	HEAD(string, interface{}, ...ReqFunc) Response
}

// TransformFunc will allow chained transformations of the Request or Response
// bodies to provide a specific value to run expectations against.
//
//		func(T<in>) (T<out>, error)
//
// `T<in>` represents any particular type. The returned `T<out>` must match the
// next transform func's argument type.
//
// The `T<out>` of the last returned transformation will be what expectations
// will be run against.
type TransformFunc interface{}

type Request interface {
	Header(string) Expectation
	Body(...TransformFunc) Expectation
	Host() Expectation

	// Request gets the actual http.Request of the underlying interface
	Request() *http.Request
}

type Response interface {
	Status() Expectation
	Header(string) Expectation
	Body(...TransformFunc) Expectation

	// Request wraps the actual request around a Request interface
	Request() Request
}

type Expectation interface {
	Equals(interface{})
	Contains(interface{})
}

type reqr struct {
	testing.TB

	handler http.Handler
}

func New(h http.Handler, t testing.TB) Reqr {
	return &reqr{
		TB: t,

		handler: h,
	}
}

func (r *reqr) makeReader(v interface{}) io.Reader {
	if v == nil {
		return nil
	}

	var reader io.Reader
	switch t := v.(type) {
	case string:
		reader = strings.NewReader(t)

	case io.Reader:
		reader = t

	default:
		w := bytes.NewBuffer(nil)
		if err := json.NewEncoder(w).Encode(v); err != nil {
			r.Fatal(err)
		}

		reader = w
	}

	return reader
}

func (r *reqr) Do(meth, path string, b interface{}, opts ...ReqFunc) Response {
	req, err := http.NewRequest(meth, path, r.makeReader(b))
	if err != nil {
		r.Fatal(err)
	}
	for _, v := range opts {
		v(req)
	}

	return r.do(req)
}

func (r *reqr) do(req *http.Request) Response {
	w := httptest.NewRecorder()

	r.handler.ServeHTTP(w, req)

	return &response{
		TB: r.TB,

		req:  req,
		resp: w,
	}
}

func (r *reqr) GET(path string, opts ...ReqFunc) Response {
	return r.Do("GET", path, nil, opts...)
}

func (r *reqr) POST(path string, b interface{}, opts ...ReqFunc) Response {
	return r.Do("POST", path, b, opts...)
}

func (r *reqr) PUT(path string, b interface{}, opts ...ReqFunc) Response {
	return r.Do("PUT", path, b, opts...)
}

func (r *reqr) PATCH(path string, b interface{}, opts ...ReqFunc) Response {
	return r.Do("PATCH", path, b, opts...)
}

func (r *reqr) DELETE(path string, b interface{}, opts ...ReqFunc) Response {
	return r.Do("DELETE", path, b, opts...)
}

func (r *reqr) OPTIONS(path string, b interface{}, opts ...ReqFunc) Response {
	return r.Do("OPTIONS", path, b, opts...)
}

func (r *reqr) HEAD(path string, b interface{}, opts ...ReqFunc) Response {
	return r.Do("HEAD", path, b, opts...)
}
