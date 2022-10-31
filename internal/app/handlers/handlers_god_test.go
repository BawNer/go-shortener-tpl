package handlers

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestMemStorage_HandlerRequest(t *testing.T) {
	test := []struct {
		name       string
		method     string
		wantStatus string
		wantBody   bool
	}{
		{
			name:       "Test POST Method",
			method:     http.MethodPost,
			wantStatus: "201 Created",
			wantBody:   true,
		},
		{
			name:       "Test GET Method",
			method:     http.MethodGet,
			wantStatus: "404 Not Found",
			wantBody:   true,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(tt.method, "http://localhost:8080/", nil)
			if err != nil {
				t.Errorf("Expect error from request!")
			}
			resp, err := http.DefaultClient.Do(request)
			if err != nil {
				t.Errorf("Expect error from request!")
			}
			defer resp.Body.Close()
			payload, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("Expect error from read body!")
			}
			if resp.Status != tt.wantStatus {
				t.Errorf("Await status %s, but get status %s", tt.wantStatus, resp.Status)
			}
			if tt.wantBody && resp.Body == nil {
				t.Errorf("Body is empty, but await body")
			}
			fmt.Println(payload)
		})
	}
}
