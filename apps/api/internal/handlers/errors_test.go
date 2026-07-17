package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/foundryhq/foundryhq/apps/api/internal/apperrors"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func recordHandleError(err error) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	handleError(c, err)
	return w
}

func TestHandleError_AppError(t *testing.T) {
	w := recordHandleError(apperrors.NotFound("task not found"))

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNotFound)
	}

	var body struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.Error.Code != string(apperrors.CodeNotFound) {
		t.Errorf("error.code = %q, want %q", body.Error.Code, apperrors.CodeNotFound)
	}
	if body.Error.Message != "task not found" {
		t.Errorf("error.message = %q, want %q", body.Error.Message, "task not found")
	}
}

func TestHandleError_ValidationIncludesField(t *testing.T) {
	w := recordHandleError(apperrors.Validation("title", "title is required"))

	var body struct {
		Error struct {
			Field string `json:"field"`
		} `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.Error.Field != "title" {
		t.Errorf("error.field = %q, want %q", body.Error.Field, "title")
	}
}

func TestHandleError_UnknownErrorIsGenericInternal(t *testing.T) {
	w := recordHandleError(errors.New("pq: relation \"tasks\" does not exist"))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}

	var body struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.Error.Code != string(apperrors.CodeInternal) {
		t.Errorf("error.code = %q, want %q", body.Error.Code, apperrors.CodeInternal)
	}
	if body.Error.Message != "internal server error" {
		t.Errorf("error.message = %q leaked internal detail, want generic message", body.Error.Message)
	}
}
