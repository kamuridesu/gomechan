// Package response provides a ResponseWriter to manage responses to be sent with http.ResponseWriter, logging information of http.Request.
//
// It keeps track of how long a request took to complete, the path, method, etc using log/slog to show output.
//
// The ResponseWriter can be used to send JSON responses, HTML responses, and static files.
//
// You can ignore logging for specific requests, useful for static files or health checks.
//
// You can also chain methods to build and send responses in one line.
//
// . . .
//
//	writer := response.New(&w, r)
//	_404 := templates.LoadHTML("404.tmpl", map[string]any{"message": "Test"})
//	writer.SetHeaders(map[string]string{"content-type": "text/html"}).Build(http.StatusOK, _404).Send()
//
// . . .
package response

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// ResponseWriter to manage response to be sent with http.ResponseWriter, logging information of http.Request.
// It keeps track of how long a request took to complete, the path, method, etc using log/slog to show output.
type ResponseWriter struct {
	status    int
	body      string
	headers   map[string]string
	start     time.Time
	ignoreLog bool
	w         *http.ResponseWriter
	r         *http.Request
}

// Creates new ResponseWriter to manage Responses.
func New(w *http.ResponseWriter, r *http.Request) ResponseWriter {
	return ResponseWriter{
		start:     time.Now(),
		ignoreLog: false,
		status:    http.StatusOK,
		headers:   map[string]string{},
		body:      "",
		w:         w,
		r:         r,
	}
}

// Builds the response.
//
// Usage:
//
//	responseWriter.Build(
//		http.StatusOK,
//		"Hello World
//	)
func (r *ResponseWriter) Build(status int, body string) *ResponseWriter {
	r.body = body
	r.status = status
	return r
}

// Sets response Headers
func (r *ResponseWriter) SetHeaders(headers map[string]string) *ResponseWriter {
	r.headers = headers
	return r
}

// Sends the response to the request, adding headers and logging the request using log/slog.
func (r *ResponseWriter) Send() error {
	for k, v := range r.headers {
		(*r.w).Header().Add(k, v)
	}
	(*r.w).WriteHeader(r.status)
	_, err := io.WriteString((*r.w), r.body)
	if err != nil {
		return err
	}
	requestTime := time.Since(r.start)
	if !r.ignoreLog {
		slog.Info(fmt.Sprintf("| %-3d | %-30v | %-15s | %-6s | %-30s",
			r.status, requestTime, strings.Split(r.r.RemoteAddr, ":")[0], r.r.Method, r.r.URL))
	}
	return nil
}

// Builds the response using map[string]any as JSON, settings content-type to application/json and sends the response
//
// Usage:
//
//	responseWriter.AsJson(
//		http.StatusOK,
//		map[string]any{"message": "Hello World"}
//	)
func (r *ResponseWriter) AsJson(status int, body map[string]any) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	r.headers["content-type"] = "application/json"
	r.Build(status, string(b))
	r.Send()
	return nil
}

// Don't log request, useful for static files or health checks
func (r *ResponseWriter) IgnoreLog() *ResponseWriter {
	r.ignoreLog = true
	return r
}
