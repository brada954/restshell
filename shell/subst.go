///////////////////////////////////////////////////////////////////////////
//
//  Substitution functions
//
//  Registered substation functions can be used in variable substitution to
//  calculate values, provide variable formatting or generate unique data
//  values. Substitution functions have the ability to key any function call
//  so the same value can be returned in a subsequent call. A different key or
//  absence of key can result in a different value.
//
//  For example, a substitution function can generate a unique ID to be used
//  in a variable or HTTP body. Using a key allows the same guid to be substituted
//  multiple times. Without a key, a function is assumed to return a different
//  value but that is not guaranteed (for example a gettime() called repeatedly
//  may return the same time due to speed of the CPU)
//
//  A package init function can register functions using RegisterSubstitutionHandler. The
//  function registered identifies a function name and function group. The function
//  name is used in the substitution process to identify the function to call. The
//  functions group membership identifies the cached data used to manage key'ed
//  instance of function data. Multiple functions in the same group can use the
//  same cache data to ensure consistency for a given key.
//
//  A function is defined as: %%funcname([key, [fmt, [option]]])%%
//  When a function is parsed, the funcname is used to identify a function to
//  call. The function is given any previous data returned from a function
//  within a group (a group shares one cache item). Groups allow multiple data elements
//  to be associated together and accessed through a single key. For example:
//
//  A function group may manage the generation and display of a random mailing address.
//  When any function in the group is called, it would generate a random mailing address
//  if one was not previously generated. When the function returns the value for substitution
//  it also returns the raw data used to generate it so the substitution package can
//  maintain that state with a key. When any other function in the group is called with
//  the same key, it will get the raw data provided.
//
//  There is a newguid() function included in the system to serve as an example.
//
package shell

import (
	"fmt"
	"regexp"
	"strings"
)

// SubstitutionHandler -- A handler returns a value for substitution given a function name. The handler may
// be given a fmt string and option string to guide the appropriate formating of data.
// The raw data is returned enabling re-use of the same data when the same key is used with a function
// in the same function group).
type SubstitutionHandler func(cache interface{}, funcName string, fmt string, option string) (value string, data interface{})

// substitutionDataCache maintains the raw data returned from functions within a function group.
// The function group name is used in the lookup
type substitutionDataCache map[string]interface{}

type substHandler struct {
	group    string
	function SubstitutionHandler
}

// Mapping of a function name to handler record identifying the group and handler
var handlerMap = make(map[string]substHandler)

var regexPattern = `%%([a-zA-Z][a-zA-Z0-9]*)\(\s*([a-zA-Z0-9_]*)\s*(?:,([a-zA-Z0-9\.]+)(?:,\s*\"([a-zA-Z0-9\.\,\;\:_\-\+\\\/\$\%\@\!\~\'\s]+?)\")?)?\s*\)%%`

// RegisterSubstitutionHandler -- Register a substitution function
func RegisterSubstitutionHandler(groupName string, funcName string, fn SubstitutionHandler) {
	if _, ok := handlerMap[funcName]; !ok {
		if IsDebugEnabled() {
			fmt.Println("Registering:", groupName, funcName)
		}
		handlerMap[funcName] = substHandler{groupName, fn}
	} else {
		panic("Duplicate substitution registration: " + groupName + "." + funcName)
	}
}

// PerformVariableSubstitution -- perform substitution on a string
func PerformVariableSubstitution(input string) string {

	var localVars = buildSubstitutionFunctionVars(input)

	var replaceStrings = make([]string, 0)

	// Filters out non-string variables
	var filter = func(k string, v interface{}) bool {
		if _, ok := v.(string); !ok {
			return false
		}
		return true
	}

	// Construct the array of strings used in variable substitution
	var replaceBuilder = func(kStr string, v interface{}) {
		if rStr, ok := v.(string); ok {
			replaceStrings = append(replaceStrings, "%%"+kStr+"%%", rStr)
		}
	}

	// Build the replacement strings from global variables
	EnumerateGlobals(replaceBuilder, filter)

	// Add the strings from the substitution function
	for k, v := range localVars {
		if IsCmdDebugEnabled() {
			fmt.Println("Adding Substitution Var: ", "%%"+k+"%% =", v)
		}
		replaceStrings = append(replaceStrings, "%%"+k+"%%", v)
	}

	// Replace all tokens in the input string
	r := strings.NewReplacer(replaceStrings...)
	return r.Replace(input)
}

// IsVariableSubstitutionComplete -- Validate that variable substitution was
// complete (no variable syntax found)
func IsVariableSubstitutionComplete(input string) bool {

	if regx, err := regexp.Compile(`\%\%.*\%\%`); err == nil {
		if regx.MatchString(input) == false {
			return true
		}
	}
	return false // Note: this is returned in error situations as well (requires investigation)
}

func buildSubstitutionFunctionVars(input string) map[string]string {
	var cache = make(substitutionDataCache, 0)
	var localVars = make(map[string]string, 0)

	pattern, _ := regexp.Compile(regexPattern)
	results := pattern.FindAllStringSubmatch(input, -1)
	for i, list := range results {
		if IsCmdDebugEnabled() && IsCmdVerboseEnabled() {
			fmt.Println("Processing group list: ", i)
			for _, m := range list {
				fmt.Println(m)
			}
		}

		fn := ""
		key := ""
		format := ""
		option := ""

		varName := strings.Trim(list[0], "%%")

		if len(list) > 1 {
			fn = list[1]
		}

		if len(list) > 2 && len(list[2]) > 0 {
			key = list[2]
		}

		if len(list) > 3 && len(list[3]) > 0 {
			format = list[3]
		}

		if len(list) > 4 && len(list[4]) > 0 {
			option = list[4]
		}

		if r, ok := handlerMap[fn]; ok {
			cachekey := r.group + "__" + key
			data, precached := cache[cachekey]
			if !precached {
				data = nil
			}

			v, c := r.function(data, fn, format, option)
			localVars[varName] = v
			cache[cachekey] = c
		}
	}
	return localVars
}
