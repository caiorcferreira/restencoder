package restencoder

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRespond(t *testing.T) {
	body := map[string]string{"msg": "Hello, world!"}

	t.Run("no options", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		Respond(recorder)

		require.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("status code", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		Respond(recorder, StatusCode(http.StatusCreated))

		require.Equal(t, http.StatusCreated, recorder.Code)
	})

	t.Run("header", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		Respond(recorder, Header("X-Header", "value"))

		require.Equal(t, "value", recorder.Header().Get("X-Header"))
	})

	t.Run("json body", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		Respond(recorder, JSONBody(body))

		require.JSONEq(t, `{"msg": "Hello, world!"}`, recorder.Body.String())
	})

	t.Run("all options", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		Respond(
			recorder,
			StatusCode(http.StatusCreated),
			Header("X-Header", "value"),
			JSONBody(body),
		)

		require.Equal(t, http.StatusCreated, recorder.Code)
		require.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("Content-Type"))
		require.Equal(t, "value", recorder.Header().Get("X-Header"))
		require.JSONEq(t, `{"msg": "Hello, world!"}`, recorder.Body.String())
	})

	t.Run("error response", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		Respond(recorder, Error(fmt.Errorf("something wrong")))

		require.Equal(t, http.StatusInternalServerError, recorder.Code)
		require.JSONEq(t, `{"error": "something wrong"}`, recorder.Body.String())
	})

	t.Run("error code", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		Respond(recorder, ErrorCode("invalid_field"))

		require.Equal(t, http.StatusInternalServerError, recorder.Code)
		require.JSONEq(t, `{"code": "invalid_field"}`, recorder.Body.String())
	})

	t.Run("error message", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		Respond(recorder, ErrorMessage("something wrong"))

		require.Equal(t, http.StatusInternalServerError, recorder.Code)
		require.JSONEq(t, `{"error": "something wrong"}`, recorder.Body.String())
	})

	t.Run("complete error response", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		Respond(recorder,
			StatusCode(http.StatusBadRequest),
			ErrorCode("invalid_field"),
			ErrorMessage("something wrong"),
		)

		require.Equal(t, http.StatusBadRequest, recorder.Code)
		require.JSONEq(t, `{"code": "invalid_field","error": "something wrong"}`, recorder.Body.String())
	})

	t.Run("should fail with invalid body", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		Respond(recorder, StatusCode(http.StatusInternalServerError), JSONBody(make(chan int))) // Providing an invalid body value

		require.Equal(t, http.StatusInternalServerError, recorder.Code)
		require.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("Content-Type"))
		require.Empty(t, recorder.Body.String()) // No response body should be written due to encoding error
	})
}
