package helpers

import (
	"fmt"
	"log"
)

type Error struct {
	Msg    string
	Reason string
}

func (e Error) Error() string {
	return fmt.Sprintf("Error: %s, reason: %s", e.Msg, e.Reason)
}

func PanicOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type HttpError struct {
	Msg    string
	Reason string
	Code   int
}

func (e HttpError) Error() string {
	return fmt.Sprintf("Error: %s, reason: %s, code: %d", e.Msg, e.Reason, e.Code)
}
