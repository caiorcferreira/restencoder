package restencoder

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResponseConfig struct {
	statusCode int
	headers    map[string]string
	body       any
}

type ResponseOption func(opts *ResponseConfig)

func StatusCode(code int) ResponseOption {
	return func(opts *ResponseConfig) {
		opts.statusCode = code
	}
}

func Header(key, value string) ResponseOption {
	return func(opts *ResponseConfig) {
		opts.headers[key] = value
	}
}

func JSONBody(b any) ResponseOption {
	return func(opts *ResponseConfig) {
		opts.headers["Content-Type"] = "application/json; charset=utf-8"
		opts.body = b
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

	if responseConfig.body != nil {
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(&responseConfig.body); err != nil {
			log.Printf("failed to encode response body: %s\n", err)

			return
		}
	}
}

type ResponseErrorConfig struct {
	StatusCode int
	ErrorCode  string
	ErrorMsg   string
}

type ResponseErrorOption func(opts *ResponseErrorConfig)

func ErrorStatusCode(code int) ResponseErrorOption {
	return func(opts *ResponseErrorConfig) {
		opts.StatusCode = code
	}
}

func ErrorCode(code string) ResponseErrorOption {
	return func(opts *ResponseErrorConfig) {
		opts.ErrorCode = code
	}
}

func ErrorStringMessage(msg string) ResponseErrorOption {
	return func(opts *ResponseErrorConfig) {
		opts.ErrorMsg = msg
	}
}

func ErrorMessage(err error) ResponseErrorOption {
	return func(opts *ResponseErrorConfig) {
		opts.ErrorMsg = err.Error()
	}
}

// ErrorResponse is the form used for API responses from failures in the API.
type ErrorResponse struct {
	Code  string `json:"code,omitempty"`
	Error string `json:"error"`
}

// RespondError translate the error and write it back to the client.
func RespondError(w http.ResponseWriter, opts ...ResponseErrorOption) {
	responseConfig := ResponseErrorConfig{}
	for _, responseOption := range opts {
		responseOption(&responseConfig)
	}

	er := ErrorResponse{
		Code:  responseConfig.ErrorCode,
		Error: responseConfig.ErrorMsg,
	}

	Respond(w, StatusCode(responseConfig.StatusCode), JSONBody(er))
}
