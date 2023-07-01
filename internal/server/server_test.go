package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/vtno/zypher/internal/server"
	"github.com/vtno/zypher/internal/server/handlers"
	"github.com/vtno/zypher/internal/store"
)

type test struct {
	name           string
	method         string
	path           string
	expectedStatus int
}

func TestServer_up(t *testing.T) {
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
	s, err := server.NewServer(server.WithPort(8081))
	if err != nil {
		t.Errorf("error creating server: %v", err)
	}
	go s.Start()
	defer s.Stop(ctx)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("http://localhost:8081%s", tt.path)
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

func TestServer_key(t *testing.T) {
	ctx := context.Background()

	// create a key to test with
	store, err := store.NewBBoltStore("zypher.db")
	if err != nil {
		t.Errorf("error creating store: %v", err)
	}
	err = store.Set("twitter#prd", "somevalue")
	if err != nil {
		t.Errorf("error setting value: %v", err)
	}
	err = store.Close()
	if err != nil {
		t.Errorf("error closing store: %v", err)
	}

	s, err := server.NewServer()
	if err != nil {
		t.Errorf("error creating server: %v", err)
	}
	go s.Start()
	defer func() {
		_ = s.Stop(ctx)
		_ = os.Remove("zypher.db")
	}()

	srvUrl := fmt.Sprintf("http://localhost:8080%s", "/key")

	t.Run("/key should return the correct key by name and env on GET", func(t *testing.T) {
		params := url.Values{}
		params.Add("name", "twitter")
		params.Add("env", "prd")
		u, _ := url.ParseRequestURI(srvUrl)
		u.RawQuery = params.Encode()
		resp, err := http.Get(fmt.Sprintf("%v", u))
		if err != nil {
			t.Errorf("error sending GET request to /key: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code to be %d, got %d", http.StatusOK, resp.StatusCode)
		}
		resBody := &handlers.KeyGetResponse{}
		err = json.NewDecoder(resp.Body).Decode(resBody)
		if err != nil {
			resp.Body.Close()
			t.Errorf("error unmarshaling response body: %v", err)
		}
		if resBody.Key != "somevalue" {
			t.Errorf("expected key to be %s, got %s", "somevalue", resBody.Key)
		}
	})
	t.Run("/key should return 201 can save a key by name and env on POST", func(t *testing.T) {
		req := handlers.KeyPostRequest{
			Name: "twitter",
			Env:  "stg",
			Key:  "supersecretkey",
		}

		reqBody, err := json.Marshal(req)
		if err != nil {
			t.Errorf("error marshaling request body: %v", err)
		}
		resp, err := http.Post(srvUrl, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Errorf("error sending POST request to /key: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected status code to be %d, got %d", http.StatusOK, resp.StatusCode)
		}
		params := url.Values{}
		params.Add("name", "twitter")
		params.Add("env", "stg")
		u, _ := url.ParseRequestURI(srvUrl)
		u.RawQuery = params.Encode()
		rawResp, err := http.Get(fmt.Sprintf("%v", u))
		kgr := &handlers.KeyGetResponse{}
		err = json.NewDecoder(rawResp.Body).Decode(kgr)
		if err != nil {
			rawResp.Body.Close()
			t.Errorf("error unmarshaling response body: %v", err)
		}
		if kgr.Key != "supersecretkey" {
			t.Errorf("expected key to be %s, got %s", "supersecretkey", kgr.Key)
		}
	})
}
