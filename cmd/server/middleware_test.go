package main

import (
	"fitness_bot/internal/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)

	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	secureHeaders(next).ServeHTTP(rr, r)

	rs := rr.Result()

	tests := []struct {
		name         string
		headerName   string
		headerValue  string
		expectedCode int
	}{
		{
			name:         "Content-Security-Policy",
			headerName:   "Content-Security-Policy",
			headerValue:  "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Referrer-Policy",
			headerName:   "Referrer-Policy",
			headerValue:  "origin-when-cross-origin",
			expectedCode: http.StatusOK,
		},
		{
			name:         "X-Content-Type-Options",
			headerName:   "X-Content-Type-Options",
			headerValue:  "nosniff",
			expectedCode: http.StatusOK,
		},
		{
			name:         "X-Frame-Options",
			headerName:   "X-Frame-Options",
			headerValue:  "deny",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, rs.Header.Get(tt.headerName), tt.headerValue)
			assert.Equal(t, rs.StatusCode, tt.expectedCode)
		})
	}

}
