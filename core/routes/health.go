package routes

import (
	"net/http"

	"github.com/kamuridesu/gomechan/core/response"
)

type Check func() error

type Health struct {
	checks []Check
}

func (h *Health) HealthCheck(w http.ResponseWriter, r *http.Request) {
	rw := response.New(&w, r)
	content := []byte("{\"status\": \"up\"}")
	status := http.StatusOK

	for _, check := range h.checks {
		if check() != nil {
			content = []byte("{\"status\": \"down\"}")
			status = http.StatusInternalServerError
			break
		}
	}

	rw.IgnoreLog().SetHeaders(map[string]string{
		"content-type": "application/json",
	}).Build(status, content).Send()
}
