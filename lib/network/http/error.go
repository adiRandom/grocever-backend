package http

import "fmt"

type Error struct {
	Msg    string
	Reason string
	Code   int
}

func (e Error) Error() string {
	return fmt.Sprintf("Error: %s, reason: %s, code: %d", e.Msg, e.Reason, e.Code)
}
