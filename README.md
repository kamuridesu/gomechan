# Gomechan

Go's net/http utilities


# Example

```go
package main

import (
	"net/http"

	"github.com/kamuridesu/gomechan/core/response"
	"github.com/kamuridesu/gomechan/core/templates"
)

func main() {
	templates, err := templates.LoadTemplateFolder("./templates")
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writer := response.New(&w, r)
		writer.AsJson(http.StatusOK, map[string]any{"status": "up"})
	})
	mux.HandleFunc("/test/", func(w http.ResponseWriter, r *http.Request) {
		writer := response.New(&w, r)
		_404 := templates.LoadHTML("404.tmpl", map[string]any{"message": "Test"})
		writer.BuildAndSend(http.StatusOK, _404, map[string]string{"content-type": "text/html"})
	})
	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		writer := response.New(&w, r)
		writer.Build(http.StatusNotFound, templates.LoadHTML("404.tmpl", map[string]any{"message": "Not Found"}))
		writer.Send()
	})
	http.ListenAndServe(":8080", mux)
}
```