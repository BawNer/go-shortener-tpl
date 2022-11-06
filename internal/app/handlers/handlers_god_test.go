package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
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
			wantBody:   false,
		},
		{
			name:       "Test GET Method",
			method:     http.MethodGet,
			wantStatus: 405,
			wantBody:   false,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			sh := &MemStorage{}
			request := httptest.NewRequest(tt.method, "http://localhost:8080", nil)
			w := httptest.NewRecorder()
			s := chi.NewRouter()
			s.Post("/", sh.HandlerPostRequest)
			s.Get("/{ID}", sh.HandlerGetRequest)
			s.ServeHTTP(w, request)
			res := w.Result()

			if res.StatusCode != tt.wantStatus {
				t.Errorf("Expected status code %d, got %d", tt.wantStatus, w.Code)
			}

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if tt.wantBody && string(resBody) == "" {
				t.Errorf("Expect body!!!")
			}
		})
	}
}
