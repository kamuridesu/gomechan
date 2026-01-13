package routes

import (
	"net/http"

	"github.com/kamuridesu/gomechan/core/response"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	rw := response.New(&w, r)
	rw.IgnoreLog().SetHeaders(map[string]string{
		"content-type": "application/json",
	}).Build(http.StatusOK, []byte("{\"status\": \"up\"}")).Send()
}
