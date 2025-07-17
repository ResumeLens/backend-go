package routes

import (
	"bytes"
	"net/http"
	"testing"
)

func TestEndpointsNot404(t *testing.T) {
	baseURL := "http://localhost:8080/api/v1"
	endpoints := []struct {
		method string
		path   string
	}{
		{"POST", "/signup"},
		{"POST", "/login"},
		{"GET",  "/validate-invite"},
		{"POST", "/accept-invite"},
		{"POST", "/refresh-token"},
		// JWT-protected endpoint:
		// {"POST", "/invite"}, // TODO: Add JWT for this endpoint
	}

	for _, ep := range endpoints {
		url := baseURL + ep.path
		var resp *http.Response
		var err error
		if ep.method == "GET" {
			resp, err = http.Get(url)
		} else {
			resp, err = http.Post(url, "application/json", bytes.NewBuffer([]byte(`{}`)))
		}
		if err != nil {
			t.Errorf("Failed to %s %s: %v", ep.method, ep.path, err)
			continue
		}
		if resp.StatusCode == http.StatusNotFound {
			t.Errorf("Endpoint %s %s returned 404", ep.method, ep.path)
		}
		resp.Body.Close()
	}
} 