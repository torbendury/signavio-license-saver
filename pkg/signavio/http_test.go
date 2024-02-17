package signavio

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("mock"))
	}))

	client := http.Client{}

	statusCode, responseBody, err := sendRequest(&client, "GET", server.URL, nil)
	if err != nil {
		t.Errorf("sendRequest() (without gzip) returned error: %v", err)
	}
	if statusCode != http.StatusOK {
		t.Errorf("sendRequest() (without gzip) statusCode: got %v, want %v", statusCode, http.StatusOK)
	}
	if string(responseBody) != "mock" {
		t.Errorf("sendRequest() (without gzip) responseBody: got %v, want %v", string(responseBody), "mock")
	}

	// test gzip response
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		io.WriteString(gz, "mock")
	}))

	statusCode, responseBody, err = sendRequest(&client, "GET", server.URL, nil)
	if err != nil {
		t.Errorf("sendRequest() returned error: %v", err)
	}
	if statusCode != http.StatusOK {
		t.Errorf("sendRequest() statusCode: got %v, want %v", statusCode, http.StatusOK)
	}
	if string(responseBody) != "mock" {
		t.Errorf("sendRequest() responseBody: got %v, want %v", string(responseBody), "mock")
	}
}

func TestNewRequest(t *testing.T) {
	client := http.Client{}
	req, err := newRequest(&client, "GET", "https://signavio.schorle", nil)
	if err != nil {
		t.Errorf("newRequest() returned error: %v", err)
	}
	if req.Host != "signavio.schorle" {
		t.Errorf("newRequest() Host: got %v, want %v", req.Host, "signavio.schorle")
	}
	if req.Method != "GET" {
		t.Errorf("newRequest() Method: got %v, want %v", req.Method, "GET")
	}

}
