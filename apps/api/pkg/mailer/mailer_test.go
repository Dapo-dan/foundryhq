package mailer

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// newTestResendSender points sender's underlying resend.Client at server
// instead of the real Resend API — client.BaseURL is exported by resend-go
// for exactly this, and client itself is reachable here because this file
// lives in the same package as ResendSender (mirroring auth_usecase_test.go's
// mockMailer, which covers the EmailSender interface contract from the
// usecase side; this covers the concrete Resend wiring itself: request
// shape, auth header, and error propagation).
func newTestResendSender(t *testing.T, server *httptest.Server) *ResendSender {
	t.Helper()
	sender := NewResendSender("test-api-key", "noreply@foundryhq.com")
	baseURL, err := url.Parse(server.URL + "/")
	if err != nil {
		t.Fatalf("parsing test server URL: %v", err)
	}
	sender.client.BaseURL = baseURL
	return sender
}

func TestResendSender_Send_Success(t *testing.T) {
	var gotAuth string
	var gotBody struct {
		From    string   `json:"from"`
		To      []string `json:"to"`
		Subject string   `json:"subject"`
		Html    string   `json:"html"`
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Errorf("decoding request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"test-email-id"}`))
	}))
	defer server.Close()

	sender := newTestResendSender(t, server)
	err := sender.Send(context.Background(), "user@example.com", "Reset your password", "<p>link</p>")
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	if gotAuth != "Bearer test-api-key" {
		t.Errorf("Authorization header = %q, want %q", gotAuth, "Bearer test-api-key")
	}
	if gotBody.From != "noreply@foundryhq.com" {
		t.Errorf("from = %q, want %q", gotBody.From, "noreply@foundryhq.com")
	}
	if len(gotBody.To) != 1 || gotBody.To[0] != "user@example.com" {
		t.Errorf("to = %v, want [user@example.com]", gotBody.To)
	}
	if gotBody.Subject != "Reset your password" {
		t.Errorf("subject = %q, want %q", gotBody.Subject, "Reset your password")
	}
	if gotBody.Html != "<p>link</p>" {
		t.Errorf("html = %q, want %q", gotBody.Html, "<p>link</p>")
	}
}

func TestResendSender_Send_WrapsProviderError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"internal error","name":"internal_server_error"}`))
	}))
	defer server.Close()

	sender := newTestResendSender(t, server)
	err := sender.Send(context.Background(), "user@example.com", "Subject", "body")
	if err == nil {
		t.Fatal("expected an error when the provider responds with a failure status")
	}
}
