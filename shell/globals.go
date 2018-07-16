//////////////////////////////////////////////////////////////////////
// Manage global data, options, and common access functions
//
// globalStore -- map of global variables and data structures
//
// The globalStore is used to hold variables that the set command
// can manage and command parser can leverage for variable substitution
// It also hold references to interfaces that enable higher level
// constructs
//
// Best practice for globalstore is to prefix the variable name with
// an indicator of the command or function of the variable. The "_"
// prefix works great to store default information that can be
// easily accessed through variable substitution. The enumeration of
// global store returns "_" prefixed variables last.
//
// Another best practice is to use $ prefix to indicate a variable may be considered
// temporary. There are commands to delete all variables starting with $.
//
//////////////////////////////////////////////////////////////////////
package shell

import (
	"errors"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/pborman/getopt/v2"
)

var (
	useDebug    *bool // Debug is intended to debug the flow of a command
	useNetDebug *bool // Special network level debugging to more indepth output from rest calls
	useVerbose  *bool // Verbose is intended to provide more detailed information from a command
	useSilent   *bool // Global silent mode
	displayHelp *bool
)

var globalStore map[string]interface{} = make(map[string]interface{}, 0)

func initGlobalStore() {
	globalStore = make(map[string]interface{}, 0)
}

func EnableGlobalOptions() {
	useDebug = getopt.BoolLong("debug", 'd', "Enable debug output globally")
	useNetDebug = getopt.BoolLong("netdb", 'n', "Enable newtwork client debug output globally")
	useVerbose = getopt.BoolLong("verbose", 'v', "Enable verbose output globally")
	useSilent = getopt.BoolLong("silent", 's', "Enable silent mode globally")
	displayHelp = getopt.BoolLong("help", 'h', "Display help")
}

func SetGlobal(key string, value interface{}) {
	globalStore[key] = value
}

// Only set the global if not initialized already
func InitializeGlobal(key string, value interface{}) {
	if _, ok := globalStore[key]; !ok {
		globalStore[key] = value
	}
}

func GetGlobal(key string) interface{} {
	if v, ok := globalStore[key]; !ok {
		return nil
	} else {
		return v
	}
}

func GetGlobalString(key string) string {
	v := GetGlobal(key)
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	default:
		panic("GetGlobalString on unsupported type")
	}
}

func GetGlobalStringWithFallback(key string, fallback string) string {
	v := GetGlobalString(key)
	if v == "" {
		return fallback
	}
	return v
}

func EnumerateGlobals(fn func(key string, value interface{}), filter func(string, interface{}) bool) {
	// Supports a best practice by separating "_" prefixed keys from others
	var keys []string
	var _keys []string

	// Build list of keys to be sorted
	for k, _ := range globalStore {
		if strings.HasPrefix(k, "_") {
			_keys = append(_keys, k)
		} else {
			keys = append(keys, k)
		}
	}

	keys = append(keys, _keys...)
	sort.Strings(keys)

	// Enumerate the keys and process the map values
	for _, v := range keys {
		if filter != nil {
			if !filter(v, globalStore[v]) {
				continue
			}
		}
		fn(v, globalStore[v])
	}
}

func RemoveGlobal(key string) {
	delete(globalStore, key)
}

//////////////////////////////////////////////////////////////////////////
// Helper functions to get global option conditionals
/////////////////////////////////////////////////////////////////////////

func SetDebug(val bool) {
	if useDebug != nil {
		*useDebug = val
	}
}

func SetSilent(val bool) {
	if useSilent != nil {
		*useSilent = val
	}
}

func SetVerbose(val bool) {
	if useVerbose != nil {
		*useVerbose = val
	}
}

func IsDebugEnabled() bool {
	return useDebug != nil && *useDebug
}

func IsVerboseEnabled() bool {
	return useVerbose != nil && *useVerbose
}

func IsSilentEnabled() bool {
	return useSilent != nil && *useSilent
}

func IsNetDebugEnabled() bool {
	return useNetDebug != nil && *useNetDebug
}

func IsDisplayHelpEnabled() bool {
	return displayHelp != nil && *displayHelp
}

////////////////////////////////////////
// Global output streams
////////////////////////////////////////

var savedOutput io.Writer = nil
var savedError io.Writer = nil
var currentConsole io.Writer = os.Stdout
var currentOutput io.Writer = os.Stdout
var currentError io.Writer = os.Stderr

// General output that goes to the console stderr and a log file
func ErrorWriter() io.Writer {
	return currentError
}

// General output that can go to the console stdout or to a log File
// Generally this is data, verbose output, and debug related to a
// function
func OutputWriter() io.Writer {
	return currentOutput
}

// Output that is intended only for the console; typically help information
// and general shell/command processor debug
func ConsoleWriter() io.Writer {
	return currentConsole
}

func SetOutput(o io.Writer) error {
	if savedOutput != nil {
		return errors.New("Already redirected")
	}
	savedOutput = currentOutput
	savedError = currentError
	currentOutput = o
	currentError = io.MultiWriter(o, savedError)
	return nil
}

func ResetOutput() (io.Writer, error) {
	if savedOutput != nil {
		ret := currentOutput
		currentOutput = savedOutput
		currentError = savedError
		savedOutput = nil
		savedError = nil
		return ret, nil
	}
	return nil, errors.New("Already reset output stream")
}
