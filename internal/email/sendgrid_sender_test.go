package email

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSendGridSender_SendPasswordReset_success(t *testing.T) {
	var gotAuth string
	var gotBody []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		var err error
		gotBody, err = io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read body: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	t.Cleanup(srv.Close)

	s := NewSendGridSender("sk-test", "noreply@example.com", srv.URL).(*SendGridSender)
	s.HTTPClient = srv.Client()

	err := s.SendPasswordReset(context.Background(), "user@example.com", "temp-secret")
	if err != nil {
		t.Fatalf("SendPasswordReset: %v", err)
	}
	if gotAuth != "Bearer sk-test" {
		t.Errorf("Authorization: got %q", gotAuth)
	}
	if !strings.Contains(string(gotBody), "user@example.com") || !strings.Contains(string(gotBody), "temp-secret") {
		t.Errorf("payload should include recipient and temp password; got %s", gotBody)
	}
}

func TestSendGridSender_SendPasswordReset_errorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	}))
	t.Cleanup(srv.Close)

	s := NewSendGridSender("sk-test", "noreply@example.com", srv.URL).(*SendGridSender)
	s.HTTPClient = srv.Client()

	err := s.SendPasswordReset(context.Background(), "user@example.com", "x")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "400") {
		t.Errorf("expected status in error, got %v", err)
	}
}
