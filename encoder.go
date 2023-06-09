package restencoder

import (
	"encoding/json"
	"log"
	"net/http"
)

// ResponseConfig holds all settings to write a REST response
type ResponseConfig struct {
	StatusCode int
	Headers    map[string]string
	JSONBody   any
}

// ResponseOption defines how to set up a response
type ResponseOption func(opts *ResponseConfig)

// StatusCode sets the response HTTP status code
func StatusCode(code int) ResponseOption {
	return func(opts *ResponseConfig) {
		opts.StatusCode = code
	}
}

// Header sets a header entry in the HTTP response
func Header(key, value string) ResponseOption {
	return func(opts *ResponseConfig) {
		opts.Headers[key] = value
	}
}

// JSONBody sets the HTTP response body to be serialized as json
func JSONBody(b any) ResponseOption {
	return func(opts *ResponseConfig) {
		opts.Headers["Content-Type"] = "application/json; charset=utf-8"
		opts.JSONBody = b
	}
}

// ErrorResponse is the contract used for failures.
type ErrorResponse struct {
	Code         string `json:"code,omitempty"`
	ErrorMessage string `json:"error,omitempty"`
}

// Error sets the error message on the body using the error content.
// It changes the response body to ErrorResponse and default the status code
// to 500 (Internal Server Error) if no error status code is set.
func Error(err error) ResponseOption {
	return func(opts *ResponseConfig) {
		errRes, ok := opts.JSONBody.(ErrorResponse)
		if !ok {
			errRes = ErrorResponse{}
		}

		errRes.ErrorMessage = err.Error()

		opts.JSONBody = errRes
		opts.Headers["Content-Type"] = "application/json; charset=utf-8"

		if opts.StatusCode >= 200 && opts.StatusCode <= 399 {
			opts.StatusCode = http.StatusInternalServerError
		}
	}
}

// ErrorCode sets the error code on the body.
// It changes the response body to ErrorResponse and default the status code
// to 500 (Internal Server Error) if no error status code is set.
func ErrorCode(errCode string) ResponseOption {
	return func(opts *ResponseConfig) {
		errRes, ok := opts.JSONBody.(ErrorResponse)
		if !ok {
			errRes = ErrorResponse{}
		}

		errRes.Code = errCode

		opts.JSONBody = errRes
		opts.Headers["Content-Type"] = "application/json; charset=utf-8"

		if opts.StatusCode >= 200 && opts.StatusCode <= 399 {
			opts.StatusCode = http.StatusInternalServerError
		}
	}
}

// ErrorMessage sets the error message on the body with the given string.
// It changes the response body to ErrorResponse and default the status code
// to 500 (Internal Server Error) if no error status code is set.
func ErrorMessage(msg string) ResponseOption {
	return func(opts *ResponseConfig) {
		errRes, ok := opts.JSONBody.(ErrorResponse)
		if !ok {
			errRes = ErrorResponse{}
		}

		errRes.ErrorMessage = msg

		opts.JSONBody = errRes
		opts.Headers["Content-Type"] = "application/json; charset=utf-8"

		if opts.StatusCode >= 200 && opts.StatusCode <= 399 {
			opts.StatusCode = http.StatusInternalServerError
		}
	}
}

func newResponseConfig() ResponseConfig {
	return ResponseConfig{
		StatusCode: http.StatusOK,
		Headers:    make(map[string]string),
	}
}

// Respond write the information back to the client.
func Respond(w http.ResponseWriter, opts ...ResponseOption) {
	responseConfig := newResponseConfig()
	for _, responseOption := range opts {
		responseOption(&responseConfig)
	}

	w.WriteHeader(responseConfig.StatusCode)

	for k, v := range responseConfig.Headers {
		w.Header().Set(k, v)
	}

	if responseConfig.JSONBody == nil {
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&responseConfig.JSONBody); err != nil {
		log.Printf("failed to encode response body: %s\n", err)

		return
	}
}
