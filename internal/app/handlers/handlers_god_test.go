package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMemStorage_HandlerRequest(t *testing.T) {
	test := []struct {
		name       string
		method     string
		wantStatus int
		wantBody   bool
	}{
		{
			name:       "Test POST Method",
			method:     http.MethodPost,
			wantStatus: 201,
			wantBody:   true,
		},
		{
			name:       "Test GET Method",
			method:     http.MethodGet,
			wantStatus: 404,
			wantBody:   true,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			sh := &MemStorage{}
			request := httptest.NewRequest(tt.method, "http://localhost:8080", nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(sh.HandlerRequest)
			h.ServeHTTP(w, request)
			res := w.Result()

			if res.StatusCode != tt.wantStatus {
				t.Errorf("Expected status code %d, got %d", tt.wantStatus, w.Code)
			}

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if !tt.wantBody && resBody != nil {
				t.Errorf("Expect body!!!")
			}
		})
	}
}
