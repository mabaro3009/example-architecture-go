package httpx

import (
	"encoding/json"
	"net/http"
)

func WriteJSONResponse(w http.ResponseWriter, statusCode int, body interface{}) error {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	buff, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = w.Write(buff)
	return err
}
