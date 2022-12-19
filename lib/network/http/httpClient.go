package http

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func GetSync[TResult any](url string) (*TResult, error) {
	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var parsed TResult
	jsonErr := json.Unmarshal(body, &parsed)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return &parsed, nil
}

func PostSync[TResult any](url string, body interface{}) (*TResult, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	res, err := http.Post(url, "application/json", io.NopCloser(bytes.NewReader(jsonBody)))
	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	resBody, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var parsed TResult
	jsonErr := json.Unmarshal(resBody, &parsed)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return &parsed, nil
}
