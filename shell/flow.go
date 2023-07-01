package shell

// FlowControl - Special Command interfaces to control execution within the command processor
type FlowControl interface {
	RequestQuit() bool
	RequestNoStep() bool
}

func DoesCommandRequestQuit(cmd interface{}) bool {
	if flow, ok := cmd.(FlowControl); ok && flow.RequestQuit() {
		return true
	}
	return false
}

func DoesCommandRequestNoStep(cmd interface{}) bool {
	if flow, ok := cmd.(FlowControl); ok && flow.RequestNoStep() {
		return true
	}
	if track, trackable := cmd.(Trackable); trackable && track.DoNotCount() {
		return true
	}
	return false
}

// FlowErrorCmd -- rename FlowAction - Special values affecting command processor
type FlowErrorCmd string

type FlowError struct {
	Message string
	Cmd     FlowErrorCmd
}

var (
	// FlowQuit - Terminate the current script without an error but do not clear LastError
	FlowQuit FlowErrorCmd = "q"
	// FlowAbort - Terminate the current script with error (typically Ctrl-C during command)
	FlowAbort FlowErrorCmd = "a"
	// FlowGo - Continue and exit single step mode
	FlowGo FlowErrorCmd = "g"
)

// NewFlowError - Return a FlowError which provides actions to cmd processor
func NewFlowError(msg string, cmd FlowErrorCmd) error {
	return FlowError{msg, cmd}
}

// Error - Return the error message for a flow error
func (f FlowError) Error() string {
	return f.Message
}

// IsFlowControl - Determines if the error as the given action associated
func IsFlowControl(err error, action FlowErrorCmd) bool {
	if f, ok := err.(FlowError); ok && f.Cmd == action {
		return true
	}
	return false
}
