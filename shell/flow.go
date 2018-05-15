package shell

// Flow Command interface to redirect flow after command execution
type FlowControl interface {
	RequestQuit() bool
	RequestNoStep() bool
}

func CommandRequestsQuit(cmd interface{}) bool {
	if flow, ok := cmd.(FlowControl); ok && flow.RequestQuit() {
		return true
	}
	return false
}

func CommandRequestsNoStep(cmd interface{}) bool {
	if flow, ok := cmd.(FlowControl); ok && flow.RequestNoStep() {
		return true
	}
	if track, trackable := cmd.(Trackable); trackable && track.DoNotCount() {
		return true
	}
	return false
}

// Flow commands affecting processor
type FlowErrorCmd string

type FlowError struct {
	Message string
	Cmd     FlowErrorCmd
}

var (
	FlowQuit  FlowErrorCmd = "q"
	FlowAbort FlowErrorCmd = "a"
	FlowGo    FlowErrorCmd = "g"
)

func NewFlowError(msg string, cmd FlowErrorCmd) error {
	return FlowError{msg, cmd}
}

func (f FlowError) Error() string {
	return f.Message
}

func IsFlowControl(err error, cmd FlowErrorCmd) bool {
	if f, ok := err.(FlowError); ok && f.Cmd == cmd {
		return true
	}
	return false
}
