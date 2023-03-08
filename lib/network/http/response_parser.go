package http

import "github.com/chebyrash/promise"

func UnwrapHttpResponse[T any](response *Response[T], err error) (*T, *Error) {
	if err != nil {
		return nil, &Error{Msg: err.Error(), Code: 500}
	}
	if response == nil {
		return nil, nil
	}
	if response.Err != "" {
		return nil, &Error{Msg: response.Err, Code: response.StatusCode}
	}
	return &response.Body, nil
}

func UnwrapHttpAsyncResponse[T any](response *promise.Promise[Response[T]]) *promise.Promise[T] {
	return promise.New[T](func(resolve func(T), reject func(error)) {
		httpResponse, err := response.Await()
		if err != nil {
			reject(&Error{Msg: err.Error(), Code: 500})
			return
		}
		if httpResponse.Err != "" {
			reject(&Error{Msg: httpResponse.Err, Code: httpResponse.StatusCode})
			return
		}
		resolve(httpResponse.Body)
	})
}
