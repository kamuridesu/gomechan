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

type ResponseWriter struct {
	status  int
	body    string
	headers map[string]string
	start   time.Time
	w       *http.ResponseWriter
	r       *http.Request
}

// Creates new ResponseWriter to manage Responses.
func New(w *http.ResponseWriter, r *http.Request) ResponseWriter {
	return ResponseWriter{
		start:   time.Now(),
		status:  http.StatusOK,
		headers: map[string]string{},
		body:    "",
		w:       w,
		r:       r,
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
func (r *ResponseWriter) Build(status int, body string) {
	r.body = body
	r.status = status
}

// Sets response Headers
func (r *ResponseWriter) SetHeaders(headers map[string]string) {
	r.headers = headers
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
	slog.Info(fmt.Sprintf("| %-3d | %-30v | %-15s | %-6s | %-30s",
		r.status, requestTime, strings.Split(r.r.RemoteAddr, ":")[0], r.r.Method, r.r.URL))
	return nil
}

// Builds the response and sends it.
//
// Usage:
//
//	responseWriter.BuildAndSend(
//		http.StatusOK,
//		"Hello World",
//		map[string]string{"content-type": "text/plain"}
//	)
func (r *ResponseWriter) BuildAndSend(status int, body string, headers map[string]string) error {
	r.SetHeaders(headers)
	r.Build(status, body)
	return r.Send()
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
