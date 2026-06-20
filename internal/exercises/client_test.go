package exercises

import (
	"fitness_bot/internal/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBuildQueryIncludesUpdatedSince(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "https://example.com/exercises", nil)
	if err != nil {
		t.Fatal(err)
	}

	upstreamReq, err := http.NewRequest(http.MethodGet, "/exercises?updated_since=2026-04-07T10:00:00Z&difficulty=beginner", nil)
	if err != nil {
		t.Fatal(err)
	}

	buildQuery(req, upstreamReq)

	query := req.URL.Query()
	assert.Equal(t, "2026-04-07T10:00:00Z", query.Get("updated_since"))
	assert.Equal(t, "beginner", query.Get("difficulty"))
}

func TestFetchExercisesReturnsReadableErrorOnUpstreamHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("<html><body>login required</body></html>"))
	}))
	defer server.Close()

	client := NewExercisesClient("token", server.URL)
	_, err := client.FetchExercises(nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "upstream exercises API returned 401 Unauthorized") {
		t.Fatalf("expected status in error, got %q", err.Error())
	}

	if !strings.Contains(err.Error(), "login required") {
		t.Fatalf("expected body preview in error, got %q", err.Error())
	}
}

func TestFetchExercisesReturnsReadableErrorOnNonJSONResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("<html><body>temporary gateway page</body></html>"))
	}))
	defer server.Close()

	client := NewExercisesClient("token", server.URL)
	_, err := client.FetchExercises(nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to decode exercises API response as JSON") {
		t.Fatalf("expected decode error context, got %q", err.Error())
	}

	if !strings.Contains(err.Error(), "response preview: <html><body>temporary gateway page</body></html>") {
		t.Fatalf("expected response preview in error, got %q", err.Error())
	}
}

func TestPreviewBodyEmptyAndTruncated(t *testing.T) {
	if got := previewBody([]byte("\n  \t  ")); got != "<empty>" {
		t.Fatalf("expected <empty>, got %q", got)
	}

	long := strings.Repeat("a", 300)
	got := previewBody([]byte(long))
	if len(got) <= 240 || !strings.HasSuffix(got, "...") {
		t.Fatalf("expected truncated preview ending with ellipsis, got len=%d preview=%q", len(got), got)
	}
}
