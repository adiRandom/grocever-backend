package http

import (
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
