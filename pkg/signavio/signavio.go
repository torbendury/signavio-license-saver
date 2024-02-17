// Signavio API client interface and related types.
package signavio

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"
)

// Job represents a scheduled user deletion job and holds the job ID.
type Job struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// JobStatus represents the status of a user deletion job. Possible values are SCHEDULED, RUNNING, COMPLETED and ERROR.
type JobStatus string

const (
	Scheduled JobStatus = "SCHEDULED"
	Running   JobStatus = "RUNNING"
	Completed JobStatus = "COMPLETED"
	Error     JobStatus = "ERROR"

	LoginAPI      = "/p/login"
	GetUsersAPI   = "/p/user?count=200&offset=0&excludeLicenses=true&timeout=900000&"
	DeleteUserAPI = "/api/v2/user-jobs/delete"
	JobStatusAPI  = "/api/v2/user-jobs/"
)

// User represents a user in the Signavio system. It holds the user's email address.
type User struct {
	Rep Rep `json:"rep"`
}

// Rep represents a user's information inside Signavio.
type Rep struct {
	Email string `json:"mail"`
}

// Client is a Signavio API client. It holds the tenant ID, the API endpoint, the API user and the API token.
type Client struct {
	TenantID   string
	Endpoint   string
	APIUser    string
	APIToken   string
	Logger     *slog.Logger
	httpClient *http.Client
}

type SignavioError struct {
	Message string
	Err     error
}

func (e *SignavioError) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.Err)
}

// New returns a new Signavio API client.
func New(tenantID string, endpoint string, apiUser string, apiToken string, logger *slog.Logger) *Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		logger.Error("can not create cookie jar", "error", err)
		os.Exit(1)
	}
	return &Client{
		TenantID: tenantID,
		Endpoint: endpoint,
		APIUser:  apiUser,
		APIToken: apiToken,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			Jar:     jar,
		},
		Logger: logger,
	}
}

// Login logs in to the Signavio API and retrieves a session token which is then stored in the client.
func (c *Client) Login() error {
	data := url.Values{
		"name":      {c.APIUser},
		"password":  {c.APIToken},
		"tokenonly": {"true"},
		"tenant":    {c.TenantID},
	}
	resp, err := c.httpClient.Post(c.Endpoint+LoginAPI, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return &SignavioError{
			Message: "can not login",
			Err:     err,
		}
	}
	if resp.StatusCode != http.StatusOK {
		return &SignavioError{
			Message: fmt.Sprintf("response status was %v", resp.StatusCode),
			Err:     errors.New("login failed"),
		}
	}
	sessionCookie := ""
	url, err := url.Parse(c.Endpoint)
	if err != nil {
		return &SignavioError{
			Message: fmt.Sprintf("can not parse url %v", c.Endpoint),
			Err:     err,
		}
	}
	for _, cookie := range c.httpClient.Jar.Cookies(url) {
		if cookie.Name == "token" {
			sessionCookie = cookie.Value
			break
		}
	}
	if sessionCookie == "" {
		return &SignavioError{
			Message: "no session cookie found",
			Err:     errors.New("login failed"),
		}
	}
	DefaultHeaders["x-signavio-id"] = sessionCookie
	return nil
}

// GetUsers retrieves all users from the Signavio system.
func (c *Client) GetUsers() (*[]User, error) {
	status, response, err := sendRequest(c.httpClient, "GET", c.Endpoint+GetUsersAPI, nil)

	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, &SignavioError{
			Message: fmt.Sprintf("response status was %v", status),
			Err:     errors.New("get users failed"),
		}
	}
	var users []User
	err = json.NewDecoder(bytes.NewReader(response)).Decode(&users)
	if err != nil {
		return nil, err
	}
	return &users, nil
}

// DeleteUser schedules the deletion of a user in the Signavio system and returns the job ID.
func (c *Client) DeleteUser(user User) (*Job, error) {
	csvData, contentType, err := csvStringToFileBuffer(fmt.Sprintf("email\n%s", user.Rep.Email))
	if err != nil {
		return nil, err
	}
	DefaultHeaders["Content-Type"] = contentType
	status, response, err := sendRequest(c.httpClient, "POST", c.Endpoint+DeleteUserAPI, csvData)
	delete(DefaultHeaders, "Content-Type")
	delete(DefaultHeaders, "Content-Disposition")
	if err != nil {
		return nil, err
	}
	if status != http.StatusCreated {
		return nil, &SignavioError{
			Message: fmt.Sprintf("response status was %v, response message was %v", status, string(response)),
			Err:     errors.New("schedule delete request failed"),
		}
	}
	var job Job
	err = json.NewDecoder(bytes.NewReader(response)).Decode(&job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// GetJobStatus retrieves the status of a user deletion job.
func (c *Client) GetJobStatus(job *Job) (*Job, error) {
	status, response, err := sendRequest(c.httpClient, "GET", c.Endpoint+JobStatusAPI+string(job.ID), nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, &SignavioError{
			Message: fmt.Sprintf("response status was %v", status),
			Err:     errors.New("get job status failed"),
		}
	}
	err = json.NewDecoder(bytes.NewReader(response)).Decode(&job)
	if err != nil {
		return nil, err
	}
	return job, nil
}
