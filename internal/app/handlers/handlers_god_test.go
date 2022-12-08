package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BawNer/go-shortener-tpl/internal/app"
	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
	"github.com/BawNer/go-shortener-tpl/internal/app/storage/file"
	"github.com/BawNer/go-shortener-tpl/internal/app/storage/memory"
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
			var repository storage.Storage
			request := httptest.NewRequest(tt.method, "http://localhost:8080", nil)
			if app.Config.FileStoragePath != "" {
				repository, _ = file.New(app.Config.FileStoragePath)
			}
			repository, _ = memory.New()

			h := NewHandler(repository)

			w := httptest.NewRecorder()
			s := chi.NewRouter()
			s.Post("/", h.HandlePostRequest)
			s.Get("/{ID}", h.HandleGetRequest)
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
