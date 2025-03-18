package client

import "net/http"

type HttpClient interface {
	// Send sends an HTTP request and returns the response body, status code, and error
	Send(req *http.Request) ([]byte, int, error)
	// Do sends an HTTP request and returns the response
	Do(req *http.Request) (*http.Response, error)
}
