package json

import (
	"encoding/json"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	type returnError struct {
		Error string `json:"error"`
	}
	errorMsg := returnError{
		Error: msg,
	}
	errResp, _ := json.Marshal(errorMsg)

	w.WriteHeader(code)
	w.Write(errResp)

}

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(v)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return err
	}
	return nil
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	resp, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}
