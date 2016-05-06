package utils

import (
	"io/ioutil"
	"net/http"

	"encoding/json"
	"fmt"

	"github.com/juju/errors"
)

// GetRequestBody reads request and returns bytes
func GetRequestBody(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return body, nil
}

// Respond sends HTTP response
func Respond(w http.ResponseWriter, data interface{}, code int) error {
	w.WriteHeader(code)

	var resp []byte
	var err error

	resp, err = json.Marshal(data) // No need HAL, if input is not valid JSON

	if err != nil {
		return errors.Trace(err)
	}

	fmt.Fprintln(w, string(resp))

	return nil
}

// UnmarshalRequest unmarshal HTTP request to given struct
func UnmarshalRequest(r *http.Request, v interface{}) error {
	body, err := GetRequestBody(r)
	if err != nil {
		return errors.Trace(err)
	}

	if err := json.Unmarshal(body, v); err != nil {
		return errors.Trace(err)
	}

	return nil
}
