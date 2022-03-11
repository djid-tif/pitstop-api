package utils

import (
	"encoding/json"
	"net/http"
	"pitstop-api/src/schemas"
)

func Prettier(w http.ResponseWriter, message string, data interface{}, status int) {
	if data == nil {
		data = struct{}{}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(schemas.Response{
		Message: message,
		Data:    data,
	})
	if err != nil {
		PrintError(err)
	}
}
