package apperrors

import (
	"errors"
	"net/http"
	"testing"
)

func TestStatusFor(t *testing.T) {
	tests := []struct {
		code Code
		want int
	}{
		{CodeValidation, http.StatusBadRequest},
		{CodeUnauthorized, http.StatusUnauthorized},
		{CodeForbidden, http.StatusForbidden},
		{CodeNotFound, http.StatusNotFound},
		{CodeConflict, http.StatusConflict},
		{CodeInternal, http.StatusInternalServerError},
		{Code("unknown_code"), http.StatusInternalServerError},
		{Code(""), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		if got := StatusFor(tt.code); got != tt.want {
			t.Errorf("StatusFor(%q) = %d, want %d", tt.code, got, tt.want)
		}
	}
}

func TestConstructors(t *testing.T) {
	if err := NotFound("task not found"); err.Code != CodeNotFound || err.Message != "task not found" {
		t.Errorf("NotFound() = %+v", err)
	}
	if err := Conflict("slug taken"); err.Code != CodeConflict {
		t.Errorf("Conflict() = %+v", err)
	}
	if err := Forbidden("not a member"); err.Code != CodeForbidden {
		t.Errorf("Forbidden() = %+v", err)
	}
	if err := Unauthorized("expired token"); err.Code != CodeUnauthorized {
		t.Errorf("Unauthorized() = %+v", err)
	}
	if err := Validation("title", "title is required"); err.Code != CodeValidation || err.Field != "title" {
		t.Errorf("Validation() = %+v", err)
	}

	cause := errors.New("connection refused")
	err := Internal(cause)
	if err.Code != CodeInternal {
		t.Errorf("Internal() Code = %v, want %v", err.Code, CodeInternal)
	}
	if err.Message != "internal server error" {
		t.Errorf("Internal() Message = %q, want a generic message", err.Message)
	}
	if !errors.Is(err, cause) {
		t.Errorf("Internal() should wrap the given cause")
	}
}

func TestErrorAs(t *testing.T) {
	var original error = NotFound("task not found")

	var target *Error
	if !errors.As(original, &target) {
		t.Fatal("errors.As should unwrap to *apperrors.Error")
	}
	if target.Code != CodeNotFound {
		t.Errorf("target.Code = %v, want %v", target.Code, CodeNotFound)
	}
}
