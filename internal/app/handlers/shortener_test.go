package handlers

import (
	"bytes"
	"encoding/json"
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

func TestMemStorage_ShortenerHandler(t *testing.T) {

	type want struct {
		status int
		body   bool
	}

	type body struct {
		URL string `json:"url"`
	}

	type args struct {
		method string
		url    string
		path   string
		body   body
		want   want
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "POST REQUEST",
			args: args{
				method: http.MethodPost,
				url:    "http://localhost:8080/api/shorten",
				path:   "/api/shorten",
				body: body{
					URL: "https://ya.ru",
				},
				want: want{
					status: 201,
					body:   true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var repository storage.Storage
			dataBody, err := json.Marshal(tt.args.body)
			if err != nil {
				t.Fatal(err)
			}

			if app.Config.FileStoragePath != "" {
				repository, _ = file.New(app.Config.FileStoragePath)
			}
			repository, _ = memory.New()

			h := NewHandler(repository)

			request := httptest.NewRequest(tt.args.method, tt.args.url, bytes.NewReader(dataBody))
			w := httptest.NewRecorder()
			s := chi.NewRouter()
			s.Post(tt.args.path, h.HandlePostRequest)
			s.ServeHTTP(w, request)
			res := w.Result()
			if res.StatusCode != tt.args.want.status {
				t.Errorf("Expected status code %d, got %d", tt.args.want.status, w.Code)
			}

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if tt.args.want.body && string(resBody) == "" {
				t.Errorf("Expect body!!!")
			}
		})
	}
}
