package util

import (
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/brada954/restshell/shell"
	"github.com/brada954/restshell/shell/modifiers"
)

type AssertCommand struct {
	exitOption      *bool
	clearOption     *bool
	newOption       *bool
	reportOption    *bool
	reportAllOption *bool
	summaryOption   *bool
	nonNilValues    *bool
	messageOption   *string
	testOption      *bool
	skipOnErrOption *bool
	expectError     *bool
	executedAsserts int
	failedAsserts   int
	totalFailures   int
	totalExecuted   int
	totalSkipped    int
	modifierOptions modifiers.ModifierOptions
	historyOptions  shell.HistoryOptions
}

// Assert -- Interface to an assertion result
type Assert interface {
	Success() bool
	Failed() bool
	Message() string
}

type assert struct {
	Passed  bool
	Text    string
	Context string
}

func (a *assert) Success() bool {
	if a == nil {
		return true
	}
	return a.Passed
}

func (a *assert) Failed() bool {
	if a == nil {
		return false
	}
	return !a.Passed
}

func (a *assert) SetContext(c string) {
	a.Context = c
}

func (a *assert) Message() string {
	if a == nil {
		return "Success inferred" // This should probably panic
	}
	if len(a.Context) > 0 {
		return a.Context + ": " + a.Text
	}
	return a.Text
}

// NewAssertFailure -- Build a failed assert
func NewAssertFailure(context string, message string, a ...interface{}) Assert {
	return &assert{false, fmt.Sprintf(message, a...), context}
}

// NewAssertSuccess -- Build a successful assert
func NewAssertSuccess(context string, message string, a ...interface{}) Assert {
	return &assert{true, fmt.Sprintf(message, a...), context}
}

// NewAssert -- builds a failed assert using non-nil error or success assert with message
func NewAssert(err error, context string, message string, a ...interface{}) Assert {
	if err == nil {
		return NewAssertSuccess(context, message, a...)
	}
	return NewAssertFailure(context, err.Error())
}

// NewAssertError -- builds a failed assert using the error
func NewAssertError(err error, context string) Assert {
	return NewAssert(err, context, "Oops, valid description missing")
}

func NewAssertCommand() *AssertCommand {
	return &AssertCommand{}
}

func (cmd *AssertCommand) GetSubCommands() []string {
	var commands = []string{"EQ", "GT", "LT", "GTE", "LTE", "NEQ", "NIL", "NNIL", "ISSTR",
		"ISINT", "ISFLOAT", "ISNUM", "ISOBJ", "ISARRAY",
		"ISDATE", "NOSTR", "NODATE", "EQDATE", "ISERR", "NOERR", "HSTATUS", "EX", "NEX", "REGMATCH"}
	return shell.SortedStringSlice(commands)
}

func (cmd *AssertCommand) AddOptions(set shell.CmdSet) {
	set.SetProgram("assert [sub command]")
	set.SetUsage(func() {
		cmd.HeaderUsage(shell.ConsoleWriter())
		set.PrintUsage(shell.ConsoleWriter())
		cmd.ExtendedUsage(shell.ConsoleWriter())
	})
	cmd.newOption = set.BoolLong("new", 'n', "Start a new set of asserts")
	cmd.exitOption = set.BoolLong("exit-onfail", 0, "Exit the command processor on failure")
	cmd.clearOption = set.BoolLong("clear", 'c', "Clear tracking of failures")
	cmd.reportOption = set.BoolLong("report", 'r', "Report failure stats for current set")
	cmd.summaryOption = set.BoolLong("report-sum", 0, "Report summary of all sets")
	cmd.reportAllOption = set.BoolLong("report-all", 'a', "Report current set and summary of all sets")
	cmd.messageOption = set.StringLong("message", 'm', "", "Display message on assert failure", "message")
	cmd.testOption = set.BoolLong("test", 0, "Use first argument as the value for testing")
	cmd.nonNilValues = set.BoolLong("non-nil", 0, "Only assert for non-nil values")
	cmd.skipOnErrOption = set.BoolLong("skip-onerr", 0, "Skip assert if tested operation failed")
	cmd.expectError = set.BoolLong("expect-error", 0, "Count failures as success")
	cmd.modifierOptions = modifiers.AddModifierOptions(set)
	cmd.historyOptions = shell.AddHistoryOptions(set, shell.AlternatePaths)
	shell.AddCommonCmdOptions(set, shell.CmdDebug, shell.CmdVerbose)
}

func (cmd *AssertCommand) HeaderUsage(w io.Writer) {
	fmt.Fprintln(w, "ASSERT COMMAND")
	fmt.Fprintln(w)
	fmt.Fprintln(w, `Assertions test values extracted from the last executed JSON request`)
	fmt.Fprintln(w, "For more information consult the repository README.md")
	fmt.Fprintln(w)
}

func (cmd *AssertCommand) ExtendedUsage(w io.Writer) {
	fmt.Fprintf(w, "\nSub Commands\n")
	lines := shell.ColumnizeTokens(cmd.GetSubCommands(), 4, 15)
	for _, v := range lines {
		fmt.Fprintf(w, "  %s\n", v)
	}
}

func (cmd *AssertCommand) Execute(args []string) error {
	*cmd.reportOption = *cmd.reportOption || *cmd.reportAllOption
	*cmd.summaryOption = *cmd.summaryOption || *cmd.reportAllOption

	if len(args) < 1 {
		if !(*cmd.clearOption || *cmd.reportOption || *cmd.summaryOption || *cmd.newOption || len(*cmd.messageOption) >= 0) {
			cmd.ExtendedUsage(shell.ConsoleWriter())
			return nil
		}
		return cmd.executeReporting()
	}

	if *cmd.clearOption || *cmd.reportOption || *cmd.summaryOption || *cmd.newOption {
		fmt.Fprintf(shell.ErrorWriter(), "Warning: --new, --clear, and --report options ignored during assert checks")
	}

	if !shell.ContainsCommand(args[0], cmd.GetSubCommands()) {
		return shell.ErrArguments
	}

	valueModifierFunc := modifiers.ConstructModifier(cmd.modifierOptions)

	result, err := shell.PeekResult(0)
	if err != nil {
		if !*cmd.testOption {
			return err
		}
	} else {
		// Conditions based on result
		if *cmd.skipOnErrOption && (result.Error != nil || result.HttpStatus >= 400) {
			cmd.totalSkipped = cmd.totalSkipped + 1
			return errors.New("WARNING: Skipping ASSERT because operation failed")
		}
	}

	var errorOccurred bool
	theAssert := cmd.executeAssertions(valueModifierFunc, result, args)
	if (*cmd.expectError && theAssert.Success()) || (!*cmd.expectError && theAssert.Failed()) {
		errorOccurred = true
		cmd.failedAsserts = cmd.failedAsserts + 1
		cmd.totalFailures = cmd.totalFailures + 1
	}
	cmd.executedAsserts = cmd.executedAsserts + 1
	cmd.totalExecuted = cmd.totalExecuted + 1

	if errorOccurred {
		err = cmd.buildErrorWithMessage(errors.New(theAssert.Message()))
		if *cmd.exitOption {
			return shell.NewFlowError(err.Error(), shell.FlowAbort)
		}
		return err
	}
	shell.OnVerbose(theAssert.Message() + "\n")
	return nil
}

func (cmd *AssertCommand) executeReporting() error {
	var reterr error

	if *cmd.reportOption {
		// Batch statistics
		if cmd.executedAsserts > 0 {
			if cmd.failedAsserts > 0 {
				fmt.Fprintf(shell.OutputWriter(),
					"Assertions failed (%d): %d out of %d assertions succeeded\n",
					cmd.failedAsserts,
					cmd.executedAsserts-cmd.failedAsserts,
					cmd.executedAsserts)
				if shell.IsCmdVerboseEnabled() && cmd.failedAsserts > 0 {
					result, err := shell.PeekResult(0)
					if err == nil {
						fmt.Fprintf(shell.OutputWriter(), "Failed Response:\n%s\n", result.Text)
					}
				}
			} else {
				fmt.Fprintf(shell.OutputWriter(), "Assertions Passed (%d)\n", cmd.executedAsserts)
			}
		}
	}

	if *cmd.summaryOption {
		// Summary causes an error to be returned and lets
		// the command processor print the error (for proper exit codes)
		if cmd.totalFailures > 0 || cmd.totalSkipped > 0 {
			reterr = errors.New(cmd.buildFailedSummaryMessage())
			reterr = cmd.buildErrorWithMessage(reterr)
			if cmd.totalFailures > 0 && *cmd.exitOption {
				reterr = shell.NewFlowError(reterr.Error(), shell.FlowAbort)
			}
		} else if cmd.totalFailures == 0 && cmd.totalExecuted > 0 {
			fmt.Fprintf(shell.OutputWriter(), "ALL ASSERTIONS PASSED (%d)\n", cmd.totalExecuted)
		}
	}

	if *cmd.newOption {
		cmd.executedAsserts = 0
		cmd.failedAsserts = 0
	}
	if *cmd.clearOption {
		cmd.executedAsserts = 0
		cmd.failedAsserts = 0
		cmd.totalFailures = 0
		cmd.totalExecuted = 0
		cmd.totalSkipped = 0
	}
	return reterr
}

func (cmd *AssertCommand) executeAssertions(valueModifierFunc modifiers.ValueModifier, result shell.Result, args []string) Assert {
	// defer func() {
	// 	if retval != nil && len(args) > 1 {
	// 		retval = retval.SetContext(args[1])
	// 	}
	// }()

	if shell.IsCmdDebugEnabled() {
		errStr := "NoErr"
		if result.Error != nil {
			errStr = "HasErr"
		}
		fmt.Fprintf(shell.ConsoleWriter(),
			"Result: %s HttpStatus:%d\n",
			errStr,
			result.HttpStatus,
		)
	}

	// Start with non-path based assertions
	if len(args) == 1 {
		switch args[0] {
		case "ISERR":
			if result.Error != nil {
				return NewAssertSuccess("Last Command", "Result was an error as expected")
			}
			return NewAssertFailure("Last Command", "Unexpectedly returned success when error expected")
		case "NOERR":
			if result.Error == nil {
				return NewAssertSuccess("Last Command", "Result was a success as expected")
			}
			return NewAssertFailure("Last Command", "Unexpectedly returned an error: "+result.Error.Error())
		default:
			return NewAssertError(shell.ErrArguments, "")
		}
	}

	if len(args) == 2 {
		switch args[0] {
		case "HSTATUS":
			if strings.ToUpper(args[1]) == "OK" || strings.ToUpper(args[1]) == "SUCCESS" {
				if result.HttpStatus == 200 || result.HttpStatus == 201 {
					return NewAssertSuccess("Last Request", "HTTP Status was %s as expected", result.HttpStatusString)
				}
			}

			status, err := strconv.Atoi(args[1])
			if err != nil {
				return NewAssertFailure("Last Request", "Invalid status argument: %s", args[1])
			}
			if result.HttpStatus != status {
				return NewAssertFailure("Last Request", "Expected status %d; got %d", status, result.HttpStatus)
			}
			return NewAssertSuccess("Last Request", "HTTP Status was %s as expected", result.HttpStatusString)
		}
	}

	// Process path based assertions
	var path = args[1]
	var node interface{}
	{
		if *cmd.testOption {
			node = args[1]
			if strings.ToLower(args[1]) == "{nil}" {
				node = nil
			}
		} else {
			var err error
			node, err = cmd.historyOptions.GetNode(path, result)
			if err != nil {
				switch args[0] {
				case "NEX":
					return NewAssertSuccess(path, "Path does not exist as expected")
				case "EX":
					return NewAssertFailure(path, "Expected path does not exist")
				default:
					return NewAssertError(err, path)
				}
			}
		}
	}

	if args[0] == "NEX" {
		return NewAssertFailure(path, "Path unexpectedly exists")
	} else if args[0] == "EX" {
		return NewAssertSuccess(path, "Path exists as expected")
	}

	// Ensure non-nil if a required condition for assert
	if node == nil && *cmd.nonNilValues {
		return nil
	}

	newnode, err := valueModifierFunc(node)
	if err != nil {
		return NewAssertError(err, "Modifiers")
	}
	node = newnode

	// If modifiers created a nil and non-nil value required, just return success
	if node == nil && *cmd.nonNilValues {
		return nil
	}

	if len(args) == 2 {
		switch args[0] {
		case "NIL":
			return NewAssert(isNil(node), path, "Node was nil as expected")
		case "NNIL":
			return NewAssert(isNotNil(node), path, "Node was not nil as expected")
		case "ISOBJ":
			return NewAssert(isObject(node), path, "Node was an object as expected")
		case "ISARRAY":
			return NewAssert(isArray(node), path, "Node was an array as expected")
		case "ISDATE":
			return NewAssert(isDate(node), path, "Node was a date as expected")
		case "NODATE":
			return NewAssert(isNotDate(node), path, "Node was not a date as expected")
		case "ISSTR":
			return NewAssert(isString(node), path, "Node was a string as expected")
		case "NOSTR":
			return NewAssert(isNotString(node), path, "Node was not a string as expected")
		case "ISINT":
			return NewAssert(isInt(node), path, "Node was an integer as expected")
		case "ISFLOAT":
			return NewAssert(isFloat(node), path, "Node  was a float as expected")
		case "ISNUM":
			if isFloat(node) == nil || isInt(node) == nil {
				return NewAssertSuccess(path, "Node was a number as expected")
			}
			return NewAssertFailure(path, "Type was not a number: %v", reflect.TypeOf(node))
		default:
			return NewAssertError(shell.ErrArguments, args[0])
		}
	}

	if len(args) == 3 {
		var value = args[2]

		switch args[0] {
		case "EQ":
			return NewAssert(isEqual(node, value), path, "Value equaled %s as expected", value)
		case "NEQ":
			return NewAssert(isNotEqual(node, value), path, "Value did not equal %s as expected", value)
		case "GT":
			return NewAssert(isGt(node, value), path, "Value was greater than %s as expected", value)
		case "GTE":
			return NewAssert(isGte(node, value), path, "Value was equal or greater than %s as expected", value)
		case "LT":
			return NewAssert(isLt(node, value), path, "Value was lessor than %s as expected", value)
		case "LTE":
			return NewAssert(isLte(node, value), path, "Value was equal or lessor than %s as expected", value)
		case "EQDATE":
			return NewAssert(isDateEqual(node, value), path, "Date was equal to %s as expected", value)
		case "REGMATCH":
			return NewAssert(isRegexMatch(node, value), path, "Value matched pattern %s as expected", value)
		default:
			return NewAssertError(shell.ErrArguments, args[0])
		}
	}

	return NewAssertError(shell.ErrArguments, args[0])
}

func (cmd *AssertCommand) buildErrorWithMessage(err error) error {
	if err == nil {
		return err
	}

	if len(*cmd.messageOption) > 0 {
		var message = err.Error()
		message = strings.Join([]string{message, *cmd.messageOption}, "\n")
		return errors.New(strings.Trim(message, "\n"))
	}
	return err
}

func (cmd *AssertCommand) buildFailedSummaryMessage() string {
	totalFailedFmt := "TOTAL ASSERTIONS FAILED: %d SKIPPED: %d EXECUTED: %d"
	return fmt.Sprintf(
		totalFailedFmt,
		cmd.totalFailures,
		cmd.totalSkipped,
		cmd.totalExecuted)
}

func isEqual(i interface{}, value string) error {
	comp, err := compare(i, value)
	if err != nil {
		return err
	}
	if comp != 0 {
		return fmt.Errorf("Values not equal: %s!=%v", value, i)
	}
	return nil
}

func isNotEqual(i interface{}, value string) error {
	comp, err := compare(i, value)
	if err != nil {
		return err
	}
	if comp == 0 {
		return fmt.Errorf("Unexpected equality: %s==%v", value, i)
	}
	return nil
}

func isGt(i interface{}, value string) error {
	comp, err := compare(i, value)
	if err != nil {
		return err
	}
	if comp > 0 {
		return nil
	}
	return fmt.Errorf("Value not greater: %s!<%v", value, i)
}

func isGte(i interface{}, value string) error {
	comp, err := compare(i, value)
	if err != nil {
		return err
	}
	if comp >= 0 {
		return nil
	}
	return fmt.Errorf("Value not greater or equal: %s!<=%v", value, i)
}

func isLt(i interface{}, value string) error {
	comp, err := compare(i, value)
	if err != nil {
		return err
	}
	if comp < 0 {
		return nil
	}
	return fmt.Errorf("Value not lessor: %s!>%v", value, i)
}

func isLte(i interface{}, value string) error {
	comp, err := compare(i, value)
	if err != nil {
		return err
	}
	if comp <= 0 {
		return nil
	}
	return fmt.Errorf("Value not lessor or equal: %s!>=%v", value, i)
}

func isNotEqualx(i interface{}, value string) error {
	switch t := i.(type) {
	case string:
		if t == value {
			return fmt.Errorf("Values unexpectedly equal: %s==%s", value, t)
		}
		return nil
	case float64:
		numValue, err := strconv.ParseFloat(value, 64)
		if shell.IsCmdDebugEnabled() {
			fmt.Printf("Debug: nodeValue: %f value: %f\n", t, numValue)
		}
		if err != nil {
			return shell.ErrInvalidValue
		}
		if math.Abs((t - numValue)) < .00001 {
			return fmt.Errorf("Value unexpectedly equal: %s==%f", value, t)
		}
		return nil
	default:
		return errors.New(shell.ErrUnexpectedType.Error() + ": " + reflect.TypeOf(i).String())
	}
}

func isRegexMatch(i interface{}, pattern string) error {
	var testValue string
	switch t := i.(type) {
	case string:
		testValue = t
	case float64:
		testValue = strconv.FormatFloat(t, 'f', -1, 64)
	case int:
		testValue = strconv.Itoa(t)
	case int64:
		testValue = strconv.FormatInt(t, 10)
	case bool:
		testValue = strconv.FormatBool(t)
	default:
		return errors.New(shell.ErrUnexpectedType.Error() + ": " + reflect.TypeOf(i).String())
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("%s: %s", pattern, err.Error())
	}
	if regex.MatchString(testValue) == false {
		return fmt.Errorf("Values does not match regex: %s!=%v", pattern, testValue)
	}
	return nil
}

func compare(i interface{}, value string) (int, error) {
	if i == nil {
		return +1, nil
	}

	switch t := i.(type) {
	case string:
		return strings.Compare(t, value), nil
	case float64:
		numValue, err := strconv.ParseFloat(value, 64)
		if shell.IsCmdDebugEnabled() {
			fmt.Fprintf(shell.ConsoleWriter(), "Debug: nodeValue: %f value: %f\n", t, numValue)
		}
		if err != nil {
			return 0, shell.ErrInvalidValue
		}
		if math.Abs((t - numValue)) < .00001 {
			return 0, nil
		} else if t > numValue {
			return +1, nil
		} else {
			return -1, nil
		}
	case int:
		numValue, err := strconv.Atoi(value)
		if err != nil {
			return 0, err
		}
		if t == numValue {
			return 0, nil
		} else if t > numValue {
			return +1, nil
		} else {
			return -1, nil
		}
	case int64:
		numValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return 0, err
		}
		if t == numValue {
			return 0, nil
		} else if t > numValue {
			return +1, nil
		} else {
			return -1, nil
		}
	case bool:
		boolval, err := strconv.ParseBool(value)
		if err != nil {
			return -1, err
		}
		if t == boolval {
			return 0, nil
		} else if boolval == false {
			return +1, nil
		} else {
			return -1, nil
		}
	default:
		return 0, errors.New(shell.ErrUnexpectedType.Error() + ": " + reflect.TypeOf(i).String())
	}
}

func isNil(i interface{}) error {
	if i == nil {
		return nil
	}
	if v, ok := i.(string); ok {
		return errors.New("Value was not nil; found value: " + v)
	}
	return errors.New("Value was not nil and not a string")
}

func isNotNil(i interface{}) error {
	if i != nil {
		return nil
	}
	return errors.New("Value was nil")
}

func isArray(i interface{}) error {
	switch i.(type) {
	case []string:
		return nil
	case []float64:
		return nil
	case []interface{}:
		return nil
	}
	return fmt.Errorf("Value was not an array: type=%v", reflect.TypeOf(i))
}

func isObject(i interface{}) error {
	switch i.(type) {
	case map[string]interface{}:
		return nil
	}
	return fmt.Errorf("Value was not an object: type=%v", reflect.TypeOf(i))
}

func isString(i interface{}) error {
	switch i.(type) {
	case string:
		return nil
	}
	return fmt.Errorf("Value was not a string: type=%v", reflect.TypeOf(i))
}

func isNotString(i interface{}) error {
	switch i.(type) {
	case string:
		return fmt.Errorf("Value was an unexpected string")
	}
	return nil
}

func isInt(i interface{}) error {
	switch i.(type) {
	case int:
		return nil
	case int32:
		return nil
	case int64:
		return nil
	}
	return fmt.Errorf("Value was not an integer: type=%v", reflect.TypeOf(i))
}

func isFloat(i interface{}) error {
	switch i.(type) {
	case float32:
		return nil
	case float64:
		return nil
	}
	return fmt.Errorf("Value was not an float: type=%v", reflect.TypeOf(i))
}

///////////////////////////////////////////////////////////////////////
// Date functions -- TODO: There are different formats to support
// Only works if expecting typical golang time displayed.
///////////////////////////////////////////////////////////////////////
func isDate(i interface{}) error {
	_, err := shell.GetValueAsDate(i)
	if err != nil {
		return err
	}
	return nil
}

func isNotDate(i interface{}) error {
	_, err := shell.GetValueAsDate(i)
	if err == nil {
		if s, ok := i.(string); ok {
			return fmt.Errorf("Value is a date: %s", s)
		} else {
			return fmt.Errorf("Value is a date: %v", i)
		}
	}
	return nil
}

func isDateEqual(i interface{}, value string) error {
	date, err := shell.GetValueAsDate(i)
	if err != nil {
		return err
	}
	expectedDate, err := shell.GetValueAsDate(value)
	if err != nil {
		return err
	}
	if date.Equal(expectedDate) {
		return nil
	}
	return fmt.Errorf("Date values not equal: %v!=%v", expectedDate, date)
}

// func onSuccessVerbose(err error, context string, format string, a ...interface{}) Assert {
// 	var assert Assert
// 	if err == nil {
// 		assert = NewAssertSuccess(context, fmt.Srintf(format, a...))
// 		if format[len(format)-1:] != "\n" {
// 			format = format + "\n"
// 		}
// 		shell.OnVerbose(format, a...)
// 	} else {
// 		assert = NewAssertFailure(context, err.Error())
// 	}
// 	return assert
// }
