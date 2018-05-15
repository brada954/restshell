package util

import (
	"os"
	"reflect"
	"testing"

	"github.com/brada954/restshell/shell"
	"github.com/pborman/getopt/v2"
)

func TestEnvSetWithFoundEnvVar(t *testing.T) {
	var trueValue = true
	expectedValue := "found_value"
	varName := "MYTEST"
	shell.RemoveGlobal(varName)

	defer func() {
		osLookupEnv = os.LookupEnv
	}()

	osLookupEnv = func(key string) (string, bool) {
		return expectedValue, true
	}

	cmd := NewSetCommand()
	cmd.AddOptions(getopt.New())
	cmd.valueIsEnvVar = &trueValue

	processArg(cmd, varName+"=xyz")

	value := shell.GetGlobalString(varName)
	if value != expectedValue {
		t.Errorf("Env variable %s!=%s", expectedValue, value)
	}
}

func TestEnvSetWithNotFoundEnvVar(t *testing.T) {
	var trueValue = true
	varName := "MYTEST2"
	shell.RemoveGlobal(varName)

	defer func() {
		osLookupEnv = os.LookupEnv
	}()

	osLookupEnv = func(key string) (string, bool) {
		return "", false
	}

	cmd := NewSetCommand()
	cmd.AddOptions(getopt.New())
	cmd.valueIsEnvVar = &trueValue

	processArg(cmd, varName+"=xyz")

	value := shell.GetGlobal(varName)
	if value != nil {
		if found, ok := value.(string); ok {
			t.Errorf("Env variable unexpectedly set to string: %s", found)
		} else {
			t.Errorf("Env variable unexpectedly set to non nil object: %s", reflect.TypeOf(value))
		}
	}
}

func TestEnvSetWithEmptyEnvVar(t *testing.T) {
	var trueValue = true
	expectedValue := ""
	varName := "MYTEST3"
	shell.RemoveGlobal(varName)

	defer func() {
		osLookupEnv = os.LookupEnv
	}()

	osLookupEnv = func(key string) (string, bool) {
		return "", true
	}

	cmd := NewSetCommand()
	cmd.AddOptions(getopt.New())
	cmd.valueIsEnvVar = &trueValue

	processArg(cmd, varName+"=xyz")

	value := shell.GetGlobalString(varName)
	if value != expectedValue {
		t.Errorf("Env variable %s!=%s", expectedValue, value)
	}
}

func TestEnvInitWithFoundEnvVar(t *testing.T) {
	var trueValue = true
	expectedValue := "original"
	foundValue := "newvalue"
	varName := "MYTEST4"
	shell.SetGlobal(varName, "original")

	defer func() {
		osLookupEnv = os.LookupEnv
	}()

	osLookupEnv = func(key string) (string, bool) {
		return foundValue, true
	}

	cmd := NewSetCommand()
	cmd.AddOptions(getopt.New())
	cmd.valueIsEnvVar = &trueValue
	cmd.initOnly = &trueValue

	processArg(cmd, varName+"=xyz")

	value := shell.GetGlobalString(varName)
	if value != expectedValue {
		t.Errorf("Env variable %s!=%s", expectedValue, value)
	}
}

func TestEnvInitWithNotFoundEnvVar(t *testing.T) {
	var trueValue = true
	varName := "MYTEST7"

	shell.RemoveGlobal(varName)

	defer func() {
		osLookupEnv = os.LookupEnv
	}()

	osLookupEnv = func(key string) (string, bool) {
		return "", false
	}

	cmd := NewSetCommand()
	cmd.AddOptions(getopt.New())
	cmd.valueIsEnvVar = &trueValue
	cmd.initOnly = &trueValue

	processArg(cmd, varName+"=xyz")

	value := shell.GetGlobal(varName)
	if value != nil {
		found, ok := value.(string)
		if ok {
			t.Errorf("Env variable unexpectedly found string: nil!=(%s)", found)
		} else {
			t.Errorf("Env variable unexpectedly found type: nil!=%v", reflect.TypeOf(value))
		}
	}
}

// Note: allow empty is ignored with env var was a value
// that leads to en empty value through modification
func TestEnvAllowEmptyWithEmptyEnvVar(t *testing.T) {
	var trueValue = true
	expectedValue := ""
	varName := "MYTEST5"

	shell.SetGlobal(varName, "initial")

	defer func() {
		osLookupEnv = os.LookupEnv
	}()

	osLookupEnv = func(key string) (string, bool) {
		return "", true
	}

	cmd := NewSetCommand()
	cmd.AddOptions(getopt.New())
	cmd.valueIsEnvVar = &trueValue
	cmd.allowEmpty = &trueValue

	processArg(cmd, varName+"=xyz")

	value := shell.GetGlobalString(varName)
	if value != expectedValue {
		t.Errorf("Env variable %s!=%s", expectedValue, value)
	}
}
