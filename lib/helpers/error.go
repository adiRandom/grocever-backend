package helpers

import "fmt"

type Error struct {
	Msg    string
	Reason string
}

func (e Error) Error() string {
	return fmt.Sprintf("Error: %s, reason: %s", e.Msg, e.Reason)
}
