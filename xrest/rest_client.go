package xrest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Mock structure used for mocking requests
type Mock struct {
	URL        string
	HTTPMethod string
	Response   *http.Response
	Err        error
}

var (
	enabledMocks = false
	mocks        = make(map[string]*Mock)
)

func getMockID(httpMethod, url string) string {
	return fmt.Sprintf("%s_%s", httpMethod, url)
}

// StartMockups enable mocking mode
func StartMockups() {
	enabledMocks = true
}

// FlushMockups clears all existing mocks from memory
func FlushMockups() {
	mocks = make(map[string]*Mock)
}

// StopMockups disable mocking mode
func StopMockups() {
	enabledMocks = false
}

// AddMock stores a new mock in memory
func AddMock(mock Mock) {
	mocks[getMockID(mock.HTTPMethod, mock.URL)] = &mock
}

// MakeRequest execute a request to a given URL with the body
func MakeRequest(method string, url string, body interface{}, headers http.Header) (*http.Response, error) {
	var jsonBytes []byte
	var err error

	if enabledMocks {
		mock := mocks[getMockID(method, url)]
		if mock != nil {
			return nil, errors.New("no mock found for given request")
		}
		return mock.Response, mock.Err
	}

	// Check if the body is already a string (Eg. JSON string)
	if w, ok := body.(string); ok {
		jsonBytes = []byte(w)
	} else {
		// Attempt to Marshal into a JSON string
		jsonBytes, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}
	request, err := http.NewRequest(method, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}
	request.Header = headers

	client := http.Client{}
	return client.Do(request)
}

// PostForm issues a POST to the specified URL, with data's keys and values URL-encoded as the request body.
func PostForm(url string, data url.Values, headers http.Header) (*http.Response, error) {
	if enabledMocks {
		mock := mocks[getMockID(http.MethodPost, url)]
		if mock != nil {
			return nil, errors.New("no mock found for given request")
		}
		return mock.Response, mock.Err
	}

	request, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	request.Header = headers

	client := http.Client{}
	return client.Do(request)
}
