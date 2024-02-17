package signavio

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	MockTenant   = "abcdefghij0123456789"
	MockEndpoint = "https://signavio.schorle"
	MockAPIUser  = "schorle@schorle.io"
	MockAPIToken = "5ch0rl3"
)

func TestNew(t *testing.T) {
	client := New(MockTenant, MockEndpoint, MockAPIUser, MockAPIToken, nil)
	if client.TenantID != MockTenant {
		t.Errorf("TenantID: got %v, want %v", client.TenantID, MockTenant)
	}
	if client.Endpoint != MockEndpoint {
		t.Errorf("Endpoint: got %v, want %v", client.Endpoint, MockEndpoint)
	}
	if client.APIUser != MockAPIUser {
		t.Errorf("ApiUser: got %v, want %v", client.APIUser, MockAPIUser)
	}
	if client.APIToken != MockAPIToken {
		t.Errorf("ApiToken: got %v, want %v", client.APIToken, MockAPIToken)
	}
}

func TestLogin(t *testing.T) {
	client := New(MockTenant, MockEndpoint, MockAPIUser, MockAPIToken, nil)
	err := client.Login()
	if err == nil {
		t.Errorf("Login() should return error: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "token", Value: "mock", Path: "/", HttpOnly: false})
		w.WriteHeader(http.StatusOK)
	}))
	client.httpClient = server.Client()
	client.httpClient.Jar, _ = cookiejar.New(nil)
	client.Endpoint = server.URL

	err = client.Login()
	if err != nil {
		t.Errorf("Login() returned error: %v", err)
	}

	url, _ := url.Parse(server.URL)
	if client.httpClient.Jar.Cookies(url)[0].Name != "token" {
		t.Errorf("Login() did not set session cookie")
	}
}
