package workouts

import (
	"fitness_bot/internal/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping)

	ts := httptest.NewTLSServer(mux)
	defer ts.Close()

	rs, err := ts.Client().Get(ts.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(body), "OK")

}
