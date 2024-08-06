package responses

import (
	"net/http"
	"net/url"
)

type Request struct {
	Method string      `json:"method"`
	Path   string      `json:"path"`
	Query  url.Values  `json:"query"`
	Body   string      `json:"body"`
	Header http.Header `json:"header"`
}
