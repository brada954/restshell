package shell

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pborman/getopt/v2"
)

type Command interface {
	Execute([]string) error
	AddOptions(set *getopt.Set)
}

type Abortable interface {
	Abort()
}

type Trackable interface {
	DoNotCount() bool
	DoNotClearError() bool
	CommandCount() int
}

type LineProcessor interface {
	ExecuteLine(line string, echoed bool) error
}

type CommandWithSubcommands interface {
	GetSubCommands() []string
}

var (
	CategoryHttp        = "Http"
	CategorySpecialized = "Specialized"
	CategoryUtilities   = "Utility"
	CategoryBenchmarks  = "Benchmark"
	CategoryTests       = "Test"
	CategoryAnalysis    = "Result Processing"
	CategoryHelp        = "Help"
)

var cmdMap map[string]Command = make(map[string]Command)
var cmdKeys map[string][]string = make(map[string][]string)
var cmdCategories []string = make([]string, 0)
var cmdSubCommands map[string][]string = make(map[string][]string)

// Cmd structures should avoid pointers to data structures so cmd structures can
// be duplicated into separate instances without data collision
func AddCommand(name string, category string, cmd Command) {
	name = strings.ToUpper(name)
	category = strings.ToLower(category)

	validateCmdEntry(name, cmd)
	ensureCategory(category)

	keys, ok := cmdKeys[category]
	if !ok {
		panic("category should exist")
	}

	cmdKeys[category] = append(keys, name)
	cmdMap[name] = cmd

	if subCmd, ok := cmd.(CommandWithSubcommands); ok {
		subcommands := subCmd.GetSubCommands()
		if len(subcommands) > 0 {
			cmdSubCommands[name] = subcommands
		}
	}
}

func ensureCategory(category string) {
	category = strings.ToLower(category)
	if keys, ok := cmdKeys[category]; !ok {
		keys = make([]string, 0)
		cmdCategories = append(cmdCategories, category)
		cmdKeys[category] = keys
	}
}

func validateCmdEntry(name string, cmd Command) {
	cmdType := reflect.TypeOf(cmd)
	for k, v := range cmdMap {
		if k == name || (v != nil && reflect.TypeOf(v) == cmdType) {
			panic("Command added more than once: " + name)
		}
	}
}

func CommandProcessesLine(cmd interface{}) bool {
	if _, isLine := cmd.(LineProcessor); isLine {
		return true
	}
	return false
}

//////////////////////////////////////////////////////////////////////
// Common commands that can be shared by commands
//////////////////////////////////////////////////////////////////////

const (
	CmdHelp int = iota
	CmdDebug
	CmdNetDebug
	CmdVerbose
	CmdTimeout
	CmdSilent
	CmdUrl
	CmdNoAuth
	CmdBasicAuth
	CmdQueryParamAuth
	CmdBenchmarks
	CmdRestclient
	CmdNoRedirect         // Encapsulated in CmdRestclient
	CmdSkipCertValidation // Encapsulated in CmdRestclient
)

const (
	OptionDefaultTimeout              = 30000 // In milliseconds
	OptionDefaultUrl                  = ""
	OptionDefaultIterations           = 10
	OptionDefaultConcurrency          = 1
	OptionDefaultIterationsThrottleMs = 0
	OptionDefaultBasicAuth            = "[user][,pwd]"
	OptionDefaultQueryParamAuth       = "[[name,value]...]"
)

type StandardOptions struct {
	helpOption           *bool
	debugOption          *bool
	netDebugOption       *bool
	verboseOption        *bool
	timeoutOption        *int64
	silentOption         *bool
	urlOption            *string
	noAuthOption         *bool
	iterationOption      *int
	concurrencyOption    *int
	throttleOption       *int
	csvOutputOption      *bool
	prettyCsvOption      *bool
	noHeaderOption       *bool
	basicAuthOption      *string
	queryParamAuthOption *string
	useLocalCertsOption  *bool
	skipCertValidation   *bool
	noRedirectOption     *bool
	reconnectOption      *bool
	warmingOption        *bool // BM warming iterations (# = concurrency)
	headersOption        *string
}

var globalOptions StandardOptions

func GetStdOptions() StandardOptions {
	return globalOptions
}

// Return a copy of the standard options with the specified
// boolean options set to true
func (o StandardOptions) Set(options ...int) StandardOptions {
	for _, v := range options {
		switch v {
		case CmdDebug:
			*o.debugOption = true
		case CmdVerbose:
			*o.verboseOption = true
		case CmdNetDebug:
			*o.netDebugOption = true
		case CmdSilent:
			*o.silentOption = true
		case CmdNoAuth:
			*o.noAuthOption = true
		case CmdSkipCertValidation:
			*o.skipCertValidation = true
		case CmdNoRedirect:
			*o.noRedirectOption = true
		default:
			panic("Illegal to set value for option:" + strconv.Itoa(v))
		}
	}
	return o
}

// Return a copy of the standard options with the specified
// options set to false for booleans or otherwise default values
func (o StandardOptions) Clear(options ...int) StandardOptions {
	for _, v := range options {
		switch v {
		case CmdDebug:
			*o.debugOption = false
		case CmdVerbose:
			*o.verboseOption = false
		case CmdNetDebug:
			*o.netDebugOption = false
		case CmdSilent:
			*o.silentOption = false
		case CmdNoAuth:
			*o.noAuthOption = false
		case CmdTimeout:
			*o.timeoutOption = OptionDefaultTimeout
		case CmdUrl:
			*o.urlOption = OptionDefaultUrl
		case CmdBenchmarks:
			*o.iterationOption = OptionDefaultIterations
			*o.throttleOption = OptionDefaultIterationsThrottleMs
			*o.concurrencyOption = OptionDefaultConcurrency
			*o.reconnectOption = false
			*o.warmingOption = false
		case CmdBasicAuth:
			*o.basicAuthOption = OptionDefaultBasicAuth
		case CmdQueryParamAuth:
			*o.queryParamAuthOption = OptionDefaultQueryParamAuth
		case CmdSkipCertValidation:
			*o.skipCertValidation = false
		case CmdNoRedirect:
			*o.noRedirectOption = false
		default:
			panic("Illegal to set value for option:" + strconv.Itoa(v))
		}
	}
	return o
}

func InitializeCommonCmdOptions(set *getopt.Set, options ...int) {
	ClearCmdOptions()
	AddCommonCmdOptions(set, options...)
}

func ClearCmdOptions() {
	globalOptions = StandardOptions{}
}

func AddCommonCmdOptions(set *getopt.Set, options ...int) {
	for _, v := range options {
		switch v {
		case CmdHelp:
			if globalOptions.helpOption == nil {
				globalOptions.helpOption = set.BoolLong("help", 'h', "Display command help")
			}
		case CmdDebug:
			if globalOptions.debugOption == nil {
				globalOptions.debugOption = set.BoolLong("debug", 'd', "Enabled debug output")
			}
		case CmdNetDebug:
			if globalOptions.netDebugOption == nil {
				globalOptions.netDebugOption = set.BoolLong("netdb", 'n', "Enabled debug output for REST client")
			}
		case CmdVerbose:
			if globalOptions.verboseOption == nil {
				globalOptions.verboseOption = set.BoolLong("verbose", 'v', "Enabled verbose output")
			}
		case CmdTimeout:
			if globalOptions.timeoutOption == nil {
				globalOptions.timeoutOption = set.Int64Long("timeout", 't', OptionDefaultTimeout, "Set the timeout for client requests in milliseconds")
			}
		case CmdSilent:
			if globalOptions.silentOption == nil {
				globalOptions.silentOption = set.BoolLong("silent", 's', "Run in silent mode")
			}
		case CmdUrl:
			if globalOptions.urlOption == nil {
				globalOptions.urlOption = set.StringLong("url", 'u', OptionDefaultUrl, "Base url for operation")
			}
		case CmdNoAuth:
			if globalOptions.noAuthOption == nil {
				globalOptions.noAuthOption = set.BoolLong("noauth", 0, "Disable auth header or context")
			}
		case CmdBasicAuth:
			if globalOptions.basicAuthOption == nil {
				globalOptions.basicAuthOption = set.StringLong("basic-auth", 0, OptionDefaultBasicAuth, "Use basic auth: "+OptionDefaultBasicAuth)
			}
		case CmdQueryParamAuth:
			if globalOptions.queryParamAuthOption == nil {
				globalOptions.queryParamAuthOption = set.StringLong("query-auth", 0, OptionDefaultQueryParamAuth, "Use query param authe: "+OptionDefaultQueryParamAuth)
			}
		case CmdBenchmarks:
			if globalOptions.iterationOption == nil {
				globalOptions.iterationOption = set.IntLong("iterations", 'i', OptionDefaultIterations, "Number of iterations for a benchmark")
			}
			if globalOptions.concurrencyOption == nil {
				globalOptions.concurrencyOption = set.IntLong("concurrency", 'c', OptionDefaultConcurrency, "Max concurrency")
			}
			if globalOptions.throttleOption == nil {
				globalOptions.throttleOption = set.IntLong("throttle", 0, OptionDefaultIterationsThrottleMs, "Delay milliseconds between iterations")
			}
			if globalOptions.csvOutputOption == nil {
				globalOptions.csvOutputOption = set.BoolLong("csv", 0, "Output results in csv format")
			}
			if globalOptions.prettyCsvOption == nil {
				globalOptions.prettyCsvOption = set.BoolLong("csv-fmt", 0, "Output formated values in CSV")
			}
			if globalOptions.noHeaderOption == nil {
				globalOptions.noHeaderOption = set.BoolLong("no-header", 0, "Do not display a header on report output")
			}
			if globalOptions.reconnectOption == nil {
				globalOptions.reconnectOption = set.BoolLong("reconnect", 0, "Force a new client and re-connect before each iteration")
			}
			if globalOptions.warmingOption == nil {
				globalOptions.warmingOption = set.BoolLong("warming", 0, "Perform warming iterations before benchmark")
			}
		case CmdRestclient:
			if globalOptions.useLocalCertsOption == nil {
				globalOptions.useLocalCertsOption = set.BoolLong("certs", 0, "Include local certs (windows may not work with system certs)")
			}
			if globalOptions.skipCertValidation == nil {
				globalOptions.skipCertValidation = set.BoolLong("nocert", 0, "Do not validate certs")
			}
			if globalOptions.noRedirectOption == nil {
				globalOptions.noRedirectOption = set.BoolLong("noredirect", 0, "Do not follow redirects on requests")
			}
			if globalOptions.headersOption == nil {
				globalOptions.headersOption = set.StringLong("headers", 0, "", "Set the headers [k=v]")
			}
		}
	}
}

func IsCmdDebugEnabled() bool {
	return IsDebugEnabled() || (globalOptions.debugOption != nil && *globalOptions.debugOption)
}

func IsCmdHelpEnabled() bool {
	return globalOptions.helpOption != nil && *globalOptions.helpOption
}

func IsCmdNetDebugEnabled() bool {
	return IsNetDebugEnabled() || (globalOptions.netDebugOption != nil && *globalOptions.netDebugOption)
}

func IsCmdVerboseEnabled() bool {
	return IsVerboseEnabled() || globalOptions.IsVerboseEnabled()
}

func GetCmdTimeoutValueMs() int64 {
	return globalOptions.GetTimeoutValueMs()
}

func IsCmdSilentEnabled() bool {
	return IsSilentEnabled() || (globalOptions.silentOption != nil && *globalOptions.silentOption)
}

func IsCmdNoAuthEnabled() bool {
	return globalOptions.IsNoAuthEnabled()
}

func IsCmdCsvFormatEnabled() bool {
	return globalOptions.IsCsvOutputEnabled()
}

func IsCmdFormattedCsvEnabled() bool {
	return globalOptions.IsFormattedCsvEnabled()
}

func IsCmdHeaderDisabled() bool {
	return globalOptions.IsHeaderDisabled()
}

func IsCmdBasicAuthEnabled() bool {
	return globalOptions.IsBasicAuthEnabled()
}

func IsCmdQueryParamAuthEnabled() bool {
	return globalOptions.IsQueryParamAuthEnabled()
}

func GetCmdUrlValue(fallback string) (result string) {
	return globalOptions.GetUrlValue(fallback)
}

func GetCmdBasicAuthContext(fallback Auth) Auth {
	return globalOptions.GetBasicAuthContext(fallback)
}

func GetCmdQueryParamAuthContext(fallback Auth) Auth {
	return globalOptions.GetQueryParamAuthContext(fallback)
}

func GetCmdHeaderValues(fallback string) string {
	return globalOptions.GetHeaderValues(fallback)
}

func GetCmdIterationValue() int {
	return globalOptions.GetCmdIterationValue()
}

func GetCmdIterationThrottleMs() int {
	return globalOptions.GetCmdIterationThrottleMs()
}

func GetCmdConcurrencyValue() int {
	return globalOptions.GetCmdConcurrencyValue()
}

func IsCmdReconnectEnabled() bool {
	return globalOptions.IsReconnectEnabled()
}

func IsCmdWarmingEnabled() bool {
	return globalOptions.IsWarmingEnabled()
}

func IsCmdLocalCertsEnabled() bool {
	return globalOptions.IsLocalCertsEnabled()
}

func IsCmdNoRedirectEnabled() bool {
	return globalOptions.IsNoRedirectEnabled()
}

func IsCmdSkipCertValidationEnabled() bool {
	return globalOptions.IsSkipCertValidationEnabled()
}

func (o *StandardOptions) IsDebugEnabled() bool {
	return o.debugOption != nil && *o.debugOption
}

func (o *StandardOptions) IsHelpEnabled() bool {
	return o.helpOption != nil && *o.helpOption
}

func (o *StandardOptions) IsNetDebugEnabled() bool {
	return o.netDebugOption != nil && *o.netDebugOption
}

func (o *StandardOptions) IsVerboseEnabled() bool {
	return o.verboseOption != nil && *o.verboseOption
}

func (o *StandardOptions) GetTimeoutValueMs() int64 {
	if o.timeoutOption != nil {
		return *o.timeoutOption
	} else {
		return OptionDefaultTimeout
	}
}

func (o *StandardOptions) IsSilentEnabled() bool {
	return o.silentOption != nil && *o.silentOption
}

func (o *StandardOptions) IsNoAuthEnabled() bool {
	return o.noAuthOption != nil && *o.noAuthOption
}

func (o *StandardOptions) IsCsvOutputEnabled() bool {
	return (o.csvOutputOption != nil && *o.csvOutputOption) || o.IsFormattedCsvEnabled()
}

func (o *StandardOptions) IsFormattedCsvEnabled() bool {
	return o.prettyCsvOption != nil && *o.prettyCsvOption
}

func (o *StandardOptions) IsHeaderDisabled() bool {
	return o.noHeaderOption != nil && *o.noHeaderOption
}

func (o *StandardOptions) IsBasicAuthEnabled() bool {
	return o.basicAuthOption != nil && *o.basicAuthOption != OptionDefaultBasicAuth
}

func (o *StandardOptions) IsQueryParamAuthEnabled() bool {
	return o.queryParamAuthOption != nil && *o.queryParamAuthOption != OptionDefaultQueryParamAuth
}

func (o *StandardOptions) GetUrlValue(fallback string) (result string) {
	defer func() {
		result = strings.Trim(result, " ") // Was trimming /; bad idea
	}()

	if o.urlOption != nil && *o.urlOption != "" {
		return *o.urlOption
	} else {
		return fallback
	}
}

func (o *StandardOptions) GetHeaderValues(fallback string) string {
	if o.headersOption != nil {
		return *o.headersOption
	} else {
		return fallback
	}
}

func (o *StandardOptions) GetBasicAuthContext(fallback Auth) Auth {
	if o.basicAuthOption != nil && *o.basicAuthOption != OptionDefaultBasicAuth {
		parts := strings.Split(*o.basicAuthOption, ",")
		var b BasicAuth
		var user string
		var pwd string
		if len(parts) == 1 {
			user = strings.TrimSpace(parts[0])
		} else if len(parts) > 1 {
			user = strings.TrimSpace(parts[0])
			pwd = strings.TrimSpace(parts[1])
		}
		b = NewBasicAuth(user, pwd)
		*o.basicAuthOption = b.UserName + "," + b.Password
		if IsCmdDebugEnabled() {
			fmt.Fprintf(ConsoleWriter(), "Returning BasicAuth for: %s\n", b.UserName)
		}
		return b
	} else {
		return fallback
	}
}

func (o *StandardOptions) GetQueryParamAuthContext(fallback Auth) Auth {
	if o.queryParamAuthOption != nil && *o.queryParamAuthOption != OptionDefaultQueryParamAuth {
		parts := strings.Split(*o.queryParamAuthOption, ",")
		var b QueryParamAuth = NewQueryParamAuth(parts...)
		if IsCmdDebugEnabled() {
			if len(b.KeyPairs) > 0 {
				key := b.KeyPairs[0].Key
				value := b.KeyPairs[0].Value
				fmt.Fprintf(ConsoleWriter(), "Returning QueryParamAuth: %s=%s\n", key, value)
			} else {
				fmt.Fprintf(ConsoleWriter(), "QueryParamAuth as no parameters\n")
			}
		}
		return b
	} else {
		return fallback
	}
}

func (o *StandardOptions) GetCmdIterationValue() int {
	if o.iterationOption != nil {
		return *o.iterationOption
	} else {
		return OptionDefaultIterations
	}
}

func (o *StandardOptions) GetCmdConcurrencyValue() int {
	if o.concurrencyOption != nil {
		return *o.concurrencyOption
	} else {
		return OptionDefaultConcurrency
	}
}

func (o *StandardOptions) GetCmdIterationThrottleMs() int {
	if o.throttleOption != nil {
		return *o.throttleOption
	} else {
		return OptionDefaultIterationsThrottleMs
	}
}

func (o *StandardOptions) IsReconnectEnabled() bool {
	return o.reconnectOption != nil && *o.reconnectOption
}

func (o *StandardOptions) IsWarmingEnabled() bool {
	return o.warmingOption != nil && *o.warmingOption
}

func (o *StandardOptions) IsLocalCertsEnabled() bool {
	return o.useLocalCertsOption != nil && *o.useLocalCertsOption
}

func (o *StandardOptions) IsSkipCertValidationEnabled() bool {
	return o.skipCertValidation != nil && *o.skipCertValidation
}

func (o *StandardOptions) IsNoRedirectEnabled() bool {
	return o.noRedirectOption != nil && *o.noRedirectOption
}

func ParseDuration(timeArg string, suffix ...string) (time.Duration, error) {
	var defaultSuffix = "ms"
	if len(suffix) > 0 {
		defaultSuffix = suffix[0]
	}

	value, err := time.ParseDuration(timeArg)
	if err != nil {
		if strings.Contains(err.Error(), "missing unit in duration") {
			value, err = time.ParseDuration(timeArg + defaultSuffix)
			if err != nil {
				return value, err
			}
		} else {
			return value, err
		}
	}
	return value, nil
}
