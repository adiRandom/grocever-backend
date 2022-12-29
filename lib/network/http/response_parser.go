package http

func ParseHttpResponse[T any](response *Response[T], err error) (*T, *Error) {
	if err != nil {
		return nil, &Error{Msg: err.Error(), Code: 500}
	}
	if response.Err != "" {
		return nil, &Error{Msg: response.Err, Code: response.StatusCode}
	}
	return &response.Body, nil
}
