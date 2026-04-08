package exercises

import (
	"net/http"
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
	if got := query.Get("updated_since"); got != "2026-04-07T10:00:00Z" {
		t.Fatalf("expected updated_since to be propagated, got %q", got)
	}

	if got := query.Get("difficulty"); got != "beginner" {
		t.Fatalf("expected difficulty to be propagated, got %q", got)
	}
}
