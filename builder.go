// Package apiclient contains some convenience methods I find myself using from
// time to time. It is not a stable collection, as even my requirements
// change from project to project, but if it ever gets there I'll update
// this description.
package apiclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// BuilderFunc is just a signature which enables to
// compose various steps of creating http requests
type BuilderFunc func(*http.Request) *http.Request

// WithHeaders sets headers on requests passed through the builder
func (bf BuilderFunc) WithHeaders(h map[string][]string) BuilderFunc {
	return func(req *http.Request) *http.Request {
		req = bf(req)
		req.Header = http.Header(h)
		return req
	}
}

// WithParams applies the params as url parameters
// It is not checked here if the resulting URL is
// of appropriate length
func (bf BuilderFunc) WithParams(params map[string][]string) BuilderFunc {
	return func(req *http.Request) *http.Request {
		// parsing a trusted source (the request's own url)
		// can get away without error check
		req = bf(req)
		values, _ := url.ParseQuery(req.URL.RawQuery)
		for key, param := range params {
			values[key] = append(values[key], param...)
		}
		req.URL.RawQuery = values.Encode()
		return req
	}
}

// WithAuth sets basic auth on the request
func (bf BuilderFunc) WithAuth(username, password string) BuilderFunc {
	return func(req *http.Request) *http.Request {
		req.SetBasicAuth(username, password)
		return bf(req)
	}
}

// MustGet ...
func MustGet(url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	return req
}

// MustPost ...
func MustPost(url string, body io.ReadCloser) *http.Request {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		panic(err)
	}
	return req
}

// MustPut ...
func MustPut(url string, body io.ReadCloser) *http.Request {
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		panic(err)
	}
	return req
}

// MustPayload takes a struct as a parameter and encodes it in the
// returned ReadCloser as a JSON object. If the passed parameter
// is not marshallable a panic is raised.
func MustPayload(v interface{}) io.ReadCloser {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return requestBody{bytes.NewBuffer(b)}
}

type requestBody struct {
	*bytes.Buffer
}

func (rb requestBody) Close() error {
	return nil
}
