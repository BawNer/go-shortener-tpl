package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
	"github.com/go-chi/chi/v5"
)

func TestMemStorage_ShortenerHandler(t *testing.T) {
	type fields struct {
		MemStorage storage.MemStorage
	}

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
		name   string
		fields fields
		args   args
	}{
		{
			name:   "POST REQUEST",
			fields: fields{},
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
			sh := &MemStorage{}
			dataBody, err := json.Marshal(tt.args.body)
			if err != nil {
				t.Fatal(err)
			}
			request := httptest.NewRequest(tt.args.method, tt.args.url, bytes.NewReader(dataBody))
			w := httptest.NewRecorder()
			s := chi.NewRouter()
			s.Post(tt.args.path, sh.ShortenerHandler)
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