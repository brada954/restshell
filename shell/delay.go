package shell

import (
	"time"
)

// Delay - Non-abortable delay function
func Delay(value time.Duration) {
	<- time.After(value)
}
