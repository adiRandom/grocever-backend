package helpers

import "lib/data/interfaces"

func SafeClose(closable interfaces.Closer) {
	err := closable.Close()
	PanicOnError(err, "Failed to close resource")
}
