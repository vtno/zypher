package server_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/vtno/zypher/internal/server"
)

func TestServer(t *testing.T) {
	type test struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}

	tests := []test{
		{
			name:           "/up should return 200 on GET",
			method:         "GET",
			path:           "/up",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "/up should return 405 on POST",
			method:         "POST",
			path:           "/up",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "/up should return 405 on PUT",
			method:         "PUT",
			path:           "/up",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "/up should return 405 on PATCH",
			method:         "PATCH",
			path:           "/up",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "/up should return 405 on DELETE",
			method:         "DELETE",
			path:           "/up",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	ctx := context.Background()
	s := server.NewServer()
	go s.Start()
	defer s.Stop(ctx)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("http://localhost:8080%s", tt.path)
			client := &http.Client{}
			req, err := http.NewRequest(tt.method, url, nil)
			if err != nil {
				t.Errorf("error creating %s request to %s: %v", tt.method, tt.path, err)
			}
			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("error sending %s request to %s: %v", tt.method, tt.path, err)
			}
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status code to be %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}
