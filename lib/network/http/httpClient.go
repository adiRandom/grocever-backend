package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chebyrash/promise"
	"io"
	"lib/helpers"
	"mime/multipart"
	"net/http"
	"os"
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
		println(readErr)
		return nil, readErr
	}

	var parsed TResult
	jsonErr := json.Unmarshal(body, &parsed)

	if jsonErr != nil {
		fmt.Printf("%v", jsonErr)
		return nil, jsonErr
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
		println(readErr)
		return nil, readErr
	}

	if len(resBody) == 0 {
		return nil, nil
	}

	var parsed TResult
	jsonErr := json.Unmarshal(resBody, &parsed)

	if jsonErr != nil {
		println(jsonErr)
		return nil, jsonErr
	}
	return &parsed, nil
}

func PostAsync[TResult any](url string, body interface{}) *promise.Promise[TResult] {
	return promise.New[TResult](func(resolve func(TResult), reject func(error)) {

		jsonBody, err := json.Marshal(body)
		if err != nil {
			reject(err)
		}

		res, err := http.Post(url, "application/json", io.NopCloser(bytes.NewReader(jsonBody)))
		if err != nil {
			reject(err)
		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		resBody, readErr := io.ReadAll(res.Body)
		if readErr != nil {
			reject(readErr)
		}

		var parsed TResult
		jsonErr := json.Unmarshal(resBody, &parsed)

		if jsonErr != nil {
			reject(jsonErr)
		}

		if err != nil {
			reject(err)
		} else {
			resolve(parsed)
		}
	})
}

type readerOrString any
type PostFormValue = readerOrString

// PostFormSync
// The values of the map must be pointer to Readers or string
func PostFormSync[TResult any](url string, values map[string]readerOrString) (*TResult, error) {
	client := &http.Client{}

	// Prepare a form that you will submit to that URL.
	var bodyBuffer bytes.Buffer
	bodyWriter := multipart.NewWriter(&bodyBuffer)

	for key, readerOrString := range values {
		if field, ok := readerOrString.(io.Closer); ok {
			defer helpers.SafeClose(field)
		}
		// Add an image file
		if field, ok := readerOrString.(*os.File); ok {
			fieldWriter, err := bodyWriter.CreateFormFile(key, field.Name())
			if err != nil {
				return nil, &helpers.Error{Msg: "Cannot create form filed"}
			}

			if _, err = io.Copy(fieldWriter, field); err != nil {
				return nil, &helpers.Error{Msg: "Cannot write the value of the field"}
			}
		} else if field, ok := readerOrString.(*multipart.File); ok {
			fieldWriter, err := bodyWriter.CreateFormFile(key, key)
			if err != nil {
				return nil, &helpers.Error{Msg: "Cannot create form filed"}
			}

			if _, err = io.Copy(fieldWriter, *field); err != nil {
				return nil, &helpers.Error{Msg: "Cannot write the value of the field"}
			}
		} else if field, ok := readerOrString.(*io.Reader); ok {
			// Add other fields
			fieldWriter, err := bodyWriter.CreateFormField(key)
			if err != nil {
				return nil, &helpers.Error{Msg: "Cannot create form filed"}
			}

			if _, err = io.Copy(fieldWriter, *field); err != nil {
				return nil, &helpers.Error{Msg: "Cannot write the value of the field"}
			}
		} else if field, ok := readerOrString.(string); ok {
			err := bodyWriter.WriteField(key, field)
			if err != nil {
				return nil, &helpers.Error{Msg: "Cannot write the value of the field"}
			}
		}

	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	helpers.SafeClose(bodyWriter)

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &bodyBuffer)
	if err != nil {
		return nil, &helpers.Error{Msg: "Cannot create the request"}
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	// Submit the request
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer helpers.SafeClose(res.Body)
	}

	resBody, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		fmt.Printf("%v", readErr)
		return nil, readErr
	}

	var parsed TResult
	jsonErr := json.Unmarshal(resBody, &parsed)

	if jsonErr != nil {
		fmt.Printf("%v", jsonErr)
		return nil, jsonErr
	}
	return &parsed, nil
}
