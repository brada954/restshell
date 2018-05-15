package shell

import (
	"time"
)

// Basic delay function (not abortable)
func Delay(value time.Duration) {
	select {
	case <-time.After(value):
	}
}
