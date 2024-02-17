package signavio

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

var (
	// DefaultHeaders are the default headers for Signavio API requests.
	DefaultHeaders = map[string]string{
		"Accept":          "application/json",
		"Accept-Encoding": "gzip, deflate, br",
		"Cache-Control":   "no-cache",
		"Charset":         "utf-8",
	}
)

func sendRequest(c *http.Client, method string, url string, body io.Reader) (int, []byte, error) {
	req, err := newRequest(c, method, url, body)
	if err != nil {
		return 0, nil, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	var respBody []byte
	if resp.Header.Get("Content-Encoding") == "gzip" {
		respBody, err = getBodyFromGzipResponse(resp)
		if err != nil {
			return 0, nil, err
		}
	} else {
		respBody, err = io.ReadAll(resp.Body)
		if err != nil {
			return 0, nil, err
		}
	}
	return resp.StatusCode, respBody, nil
}

func newRequest(r *http.Client, method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range DefaultHeaders {
		req.Header.Add(k, v)
	}
	return req, nil
}

func getBodyFromGzipResponse(resp *http.Response) ([]byte, error) {
	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func csvStringToFileBuffer(csvString string) (io.Reader, string, error) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, "data.csv"))
	header.Set("Content-Type", "text/csv")
	formWriter, err := writer.CreatePart(header)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(formWriter, strings.NewReader(csvString))
	if err != nil {
		panic(err)
	}

	if err = writer.Close(); err != nil {
		panic(err)
	}
	return &buffer, writer.FormDataContentType(), nil
}
