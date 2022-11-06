package helpers

import (
	"fmt"
	"lib/data/interfaces"
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

func SafeClose(closable interfaces.Closable) {
	err := closable.Close()
	PanicOnError(err, "Failed to close resource")
}
