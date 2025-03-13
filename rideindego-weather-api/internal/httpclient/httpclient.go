package httpclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func HTTPRequest(
	method string,
	headers map[string]string,
	endpoint string,
	payload *bytes.Buffer,
	resp any,
) (status int, e error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return 0, err
	}

	if payload == nil {
		payload = bytes.NewBuffer([]byte{})
	}

	httpReq, err := http.NewRequest(method, u.String(), payload)
	if err != nil {
		return 0, err
	}

	// append headers so as not to replace it
	for key, val := range headers {
		httpReq.Header.Add(key, val)
	}

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return 0, err
	}

	// if caller does not specify body, retrieve the status code only
	if resp == nil {
		return httpResp.StatusCode, nil
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return 0, err
	}

	if err := json.Unmarshal(body, resp); err != nil {
		return 0, errors.Join(err, fmt.Errorf("server returned %s - %v", httpResp.Status, string(body)[:64]))
	}

	return httpResp.StatusCode, nil
}
