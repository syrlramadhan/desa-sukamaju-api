package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/syrlramadhan/desa-sukamaju-api/dto"
)

// WriteJSONSuccess digunakan untuk mengirim response sukses dalam format JSON
func WriteJSONSuccess(w http.ResponseWriter, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := dto.ListResponseOK{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Data:    data,
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}

// WriteJSONNoData digunakan untuk mengirim response sukses dalam format JSON
func WriteJSONNoData(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := dto.ListResponseNoData{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}

// WriteJSONLogin digunakan untuk mengirim response sukses dalam format JSON
func WriteJSONLogin(w http.ResponseWriter, token string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := dto.LoginResponse{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Token:   token,
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}

// WriteJSONError untuk mengirim response error dengan format JSON
func WriteJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dto.ListResponseError{
		Code:    statusCode,
		Status:  http.StatusText(statusCode),
		Message: message,
	})
}
