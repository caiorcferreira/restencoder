package restencoder

import (
	"encoding/json"
	"log"
	"net/http"
)

// ResponseConfig holds all settings to write a REST response
type ResponseConfig struct {
	statusCode int
	headers    map[string]string
	body       any
}

// ResponseOption defines how to set up a response
type ResponseOption func(opts *ResponseConfig)

// StatusCode sets the response HTTP status code
func StatusCode(code int) ResponseOption {
	return func(opts *ResponseConfig) {
		opts.statusCode = code
	}
}

// Header sets a header entry in the HTTP response
func Header(key, value string) ResponseOption {
	return func(opts *ResponseConfig) {
		opts.headers[key] = value
	}
}

// JSONBody sets the HTTP response body to be serialized as json
func JSONBody(b any) ResponseOption {
	return func(opts *ResponseConfig) {
		opts.headers["Content-Type"] = "application/json; charset=utf-8"
		opts.body = b
	}
}

// ErrorResponse is the form used for API responses from failures in the API.
type ErrorResponse struct {
	Code         string `json:"code,omitempty"`
	ErrorMessage string `json:"error,omitempty"`
}

func Error(err error) ResponseOption {
	return func(opts *ResponseConfig) {
		errRes, ok := opts.body.(ErrorResponse)
		if !ok {
			errRes = ErrorResponse{}
		}

		errRes.ErrorMessage = err.Error()

		opts.body = errRes
		opts.headers["Content-Type"] = "application/json; charset=utf-8"

		if opts.statusCode >= 200 && opts.statusCode <= 399 {
			opts.statusCode = http.StatusInternalServerError
		}
	}
}

func ErrorCode(errCode string) ResponseOption {
	return func(opts *ResponseConfig) {
		errRes, ok := opts.body.(ErrorResponse)
		if !ok {
			errRes = ErrorResponse{}
		}

		errRes.Code = errCode

		opts.body = errRes
		opts.headers["Content-Type"] = "application/json; charset=utf-8"

		if opts.statusCode >= 200 && opts.statusCode <= 399 {
			opts.statusCode = http.StatusInternalServerError
		}
	}
}

func ErrorMessage(msg string) ResponseOption {
	return func(opts *ResponseConfig) {
		errRes, ok := opts.body.(ErrorResponse)
		if !ok {
			errRes = ErrorResponse{}
		}

		errRes.ErrorMessage = msg

		opts.body = errRes
		opts.headers["Content-Type"] = "application/json; charset=utf-8"

		if opts.statusCode >= 200 && opts.statusCode <= 399 {
			opts.statusCode = http.StatusInternalServerError
		}
	}
}

func newResponseConfig() ResponseConfig {
	return ResponseConfig{
		statusCode: http.StatusOK,
		headers:    make(map[string]string),
	}
}

// Respond write the information back to the client.
func Respond(w http.ResponseWriter, opts ...ResponseOption) {
	responseConfig := newResponseConfig()
	for _, responseOption := range opts {
		responseOption(&responseConfig)
	}

	w.WriteHeader(responseConfig.statusCode)

	for k, v := range responseConfig.headers {
		w.Header().Set(k, v)
	}

	if responseConfig.body == nil {
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&responseConfig.body); err != nil {
		log.Printf("failed to encode response body: %s\n", err)

		return
	}
}
