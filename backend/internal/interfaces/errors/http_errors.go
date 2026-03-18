package errors

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"asana-clone-backend/internal/domain/shared"
)

// ErrorResponse is the standard error payload returned by the API.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

// MapDomainError converts a domain error into an HTTP status code and error response.
func MapDomainError(err error) (int, ErrorResponse) {
	domainErr, ok := err.(*shared.DomainError)
	if !ok {
		log.Printf("ERROR: non-domain error: %T: %v", err, err)
		return http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "an unexpected error occurred",
		}
	}

	status := mapCodeToHTTPStatus(domainErr.Code)

	return status, ErrorResponse{
		Code:    domainErr.Code,
		Message: domainErr.Message,
		Field:   domainErr.Field,
	}
}

func mapCodeToHTTPStatus(code string) int {
	switch {
	case code == "NOT_FOUND" || strings.HasSuffix(code, "_NOT_FOUND"):
		return http.StatusNotFound
	case code == "UNAUTHORIZED" || code == "INVALID_CREDENTIALS":
		return http.StatusUnauthorized
	case code == "FORBIDDEN":
		return http.StatusForbidden
	case code == "ALREADY_EXISTS" || code == "EMAIL_TAKEN" || strings.HasPrefix(code, "ALREADY_"):
		return http.StatusConflict
	case strings.HasPrefix(code, "INVALID_") || code == "VALIDATION_ERROR":
		return http.StatusBadRequest
	case code == "SECTION_MISMATCH" || code == "SECTION_NOT_FOUND" || code == "LABEL_NOT_FOUND" || code == "USER_NOT_FOUND":
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// RespondWithError writes a JSON error response based on the given error.
func RespondWithError(w http.ResponseWriter, err error) {
	status, errResp := MapDomainError(err)
	RespondWithJSON(w, status, errResp)
}

// RespondWithJSON writes a JSON response with the given status code and payload.
func RespondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
