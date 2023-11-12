package utils

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
)

func HttpCall(context context.Context, url string, authToken string, verb string, payload []byte) (body []byte, httpstatus int, err error) {
	req, err := http.NewRequestWithContext(context, verb, url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	if authToken != "" {
		req.Header.Set("Authorization", authToken)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		//fmt.Println(err)
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return body, resp.StatusCode, nil
	}

	return body, resp.StatusCode, errors.New("server returned not success status code")
}