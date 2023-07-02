package shell

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//////////////////////////////////////////////////////////////////////
// Common options that can be shared by commands
//////////////////////////////////////////////////////////////////////

// Options to include with a command
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
	CmdFormatOutput
)

// Default values for options
const (
	OptionDefaultTimeout              = 30000 // In milliseconds
	OptionDefaultUrl                  = ""
	OptionDefaultIterations           = 0
	OptionDefaultDuration             = ""
	OptionDefaultConcurrency          = 1
	OptionDefaultIterationsThrottleMs = 0
	OptionDefaultBasicAuth            = "[user][,pwd]"
	OptionDefaultQueryParamAuth       = "[[name,value]...]"
	OptionDefaultOutputFile           = ""
)

// Structure for all common option values
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
	durationOption       *string
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
	headerOption         *StringList
	shortOutputOption    *bool
	bodyOutputOption     *bool
	headerOutputOption   *bool
	cookieOutputOption   *bool
	fullOutputOption     *bool
	fileOutputOption     *string
	prettyPrintOption    *bool
	requestOutputOption  *bool
}

// Global options populated by a command being run
var globalOptions StandardOptions

func GetStdOptions() StandardOptions {
	return globalOptions
}

// Set -- Return a copy of the standard options with the specified
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

// Clear -- Return a copy of the standard options with the specified
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
			*o.durationOption = OptionDefaultDuration
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

// InitializeCommonCmdOptions initialize common command options
func InitializeCommonCmdOptions(set CmdSet, options ...int) {
	ClearCmdOptions()
	AddCommonCmdOptions(set, options...)
}

// ClearCmdOptions -- Clear common command options by setting a new configuration.
func ClearCmdOptions() {
	globalOptions = StandardOptions{}
}

// AddCommonCmdOptions -- Add the given command options to the options supported
// by the current executing command
func AddCommonCmdOptions(set CmdSet, options ...int) {
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
				globalOptions.queryParamAuthOption = set.StringLong("query-auth", 0, OptionDefaultQueryParamAuth, "Use query parameters for auth: "+OptionDefaultQueryParamAuth)
			}
		case CmdBenchmarks:
			if globalOptions.iterationOption == nil {
				globalOptions.iterationOption = set.IntLong("iterations", 'i', OptionDefaultIterations, "Maximum iterations for a benchmark")
			}
			if globalOptions.durationOption == nil {
				globalOptions.durationOption = set.StringLong("duration", 0, OptionDefaultDuration, "Maximum duration of a benchmark")
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
				globalOptions.headersOption = set.StringLong("headers", 0, "", "Set the headers [k=v,k=v]")
			}
			if globalOptions.headerOption == nil {
				globalOptions.headerOption = set.StringListLong("header", 0, "Set a header [k=v]")
			}
		case CmdFormatOutput:
			if globalOptions.shortOutputOption == nil {
				globalOptions.shortOutputOption = set.BoolLong("out-short", 0, "Output the short response (overrides verbose)")
			}
			if globalOptions.bodyOutputOption == nil {
				globalOptions.bodyOutputOption = set.BoolLong("out-body", 0, "Output the response body")
			}
			if globalOptions.headerOutputOption == nil {
				globalOptions.headerOutputOption = set.BoolLong("out-header", 0, "Output response headers")
			}
			if globalOptions.cookieOutputOption == nil {
				globalOptions.cookieOutputOption = set.BoolLong("out-cookie", 0, "Output response cookies")
			}
			if globalOptions.fullOutputOption == nil {
				globalOptions.fullOutputOption = set.BoolLong("out-full", 0, "Output all response data")
			}
			if globalOptions.fileOutputOption == nil {
				globalOptions.fileOutputOption = set.StringLong("out-file", 0, "", "Output result to a file", "file")
			}
			if globalOptions.prettyPrintOption == nil {
				globalOptions.prettyPrintOption = set.BoolLong("pretty", 0, "Pretty print output")
			}
			if globalOptions.requestOutputOption == nil {
				globalOptions.requestOutputOption = set.BoolLong("out-request", 0, "Output the request body")
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

func GetCmdHeaderValues() []string {
	return globalOptions.GetHeaderValues()
}

func GetCmdIterationValue() int {
	return globalOptions.GetCmdIterationValue()
}

func GetCmdDurationValueWithFallback(d time.Duration) time.Duration {
	return globalOptions.GetCmdDurationValueWithFallback(0)
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

func IsCmdFullOutputEnabled() bool {
	return (globalOptions.fullOutputOption != nil && *globalOptions.fullOutputOption)
}

func IsCmdOutputHeaderEnabled() bool {
	return (globalOptions.headerOutputOption != nil && *globalOptions.headerOutputOption) || IsCmdFullOutputEnabled()
}

func IsCmdOutputCookieEnabled() bool {
	return (globalOptions.cookieOutputOption != nil && *globalOptions.cookieOutputOption) || IsCmdFullOutputEnabled()
}

func IsCmdOutputShortEnabled() bool {
	return (globalOptions.shortOutputOption != nil && *globalOptions.shortOutputOption)
}

func IsCmdOutputBodyEnabled() bool {
	return (globalOptions.bodyOutputOption != nil && *globalOptions.bodyOutputOption) || (IsCmdVerboseEnabled() && !IsCmdOutputShortEnabled())
}

func IsCmdOutputRequestEnabled() bool {
	return (globalOptions.requestOutputOption != nil && *globalOptions.requestOutputOption)
}

func GetCmdOutputFileName() string {
	return globalOptions.GetCmdOutputFileName()
}

func IsCmdPrettyPrintEnabled() bool {
	return (globalOptions.prettyPrintOption != nil && *globalOptions.prettyPrintOption)
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
	}
	return fallback
}

func (o *StandardOptions) GetHeaderValues() []string {
	result := make([]string, 0)
	if o.headersOption != nil {
		for _, v := range strings.Split(*o.headersOption, ",") {
			if len(v) > 0 {
				result = append(result, v)
			}
		}
	}
	if o.headerOption != nil {
		result = append(result, (*o.headerOption).Values...)
	}
	return result
}

// GetBasicAuthContext -- get the Auth context for the basic auth parameters specified
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

// GetQueryParamAuthContext -- return the Auth context for query parameters specified in command
func (o *StandardOptions) GetQueryParamAuthContext(fallback Auth) Auth {
	if o.queryParamAuthOption != nil && *o.queryParamAuthOption != OptionDefaultQueryParamAuth {
		parts := strings.Split(*o.queryParamAuthOption, ",")
		var auth = NewQueryParamAuth(parts...)
		if IsCmdDebugEnabled() {
			if len(auth.KeyPairs) > 0 {
				key := auth.KeyPairs[0].Key
				value := auth.KeyPairs[0].Value
				fmt.Fprintf(ConsoleWriter(), "Returning QueryParamAuth: %s=%s\n", key, value)
			} else {
				fmt.Fprintf(ConsoleWriter(), "QueryParamAuth as no parameters\n")
			}
		}
		return auth
	} else {
		return fallback
	}
}

// GetCmdIterationValue -- Get the iteration value and use the default if not set
func (o *StandardOptions) GetCmdIterationValue() int {
	if o.iterationOption != nil {
		return *o.iterationOption
	}
	return OptionDefaultIterations
}

// GetCmdIterationValueWithFallback -- Get the iteration value and use the default if not set
func (o *StandardOptions) GetCmdIterationValueWithFallback(d int) int {
	if o.iterationOption != nil && *o.iterationOption != OptionDefaultIterations {
		return *o.iterationOption
	}
	return d
}

// GetCmdDurationValueWithFallback -- Get the duration parameter use default if not set
func (o *StandardOptions) GetCmdDurationValueWithFallback(d time.Duration) time.Duration {
	if o.durationOption != nil {
		if dur, err := ParseDuration(*o.durationOption); err != nil {
			return dur
		}
	}
	return d
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

func (o *StandardOptions) IsFullOutputEnabled() bool {
	return o.fullOutputOption != nil && *o.fullOutputOption
}

func (o *StandardOptions) IsOutputHeaderEnabled() bool {
	return o.headerOutputOption != nil && *o.headerOutputOption
}

func (o *StandardOptions) IsOutputCookieEnabled() bool {
	return o.cookieOutputOption != nil && *o.cookieOutputOption
}

func (o *StandardOptions) IsOutputShortEnabled() bool {
	return o.shortOutputOption != nil && *o.shortOutputOption
}

func (o *StandardOptions) IsOutputBodyEnabled() bool {
	return o.bodyOutputOption != nil && *o.bodyOutputOption
}

func (o *StandardOptions) IsOutputRequestEnabled() bool {
	return o.requestOutputOption != nil && *o.requestOutputOption
}

func (o *StandardOptions) GetCmdOutputFileName() string {
	if o.fileOutputOption != nil && *o.fileOutputOption != OptionDefaultOutputFile {
		return *o.fileOutputOption
	}
	return OptionDefaultOutputFile
}

func (o *StandardOptions) IsPrettyPrintEnabled() bool {
	return o.prettyPrintOption != nil && *o.prettyPrintOption
}

// ParseDuration -- parses a duration value from text that may include
// time sufix [ms(default), S, see time.ParseDuration]
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
