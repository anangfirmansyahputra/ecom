package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anangfirmansyahp5/ecom/types"
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, r types.Response) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(r)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	response := types.Response{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	}

	WriteJSON(w, status, response)
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}
