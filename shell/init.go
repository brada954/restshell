package shell

import (
	"os"
	"path/filepath"
	"strings"
)

var InitDirectory string = ""
var ExecutableDirectory string = ""
var initialized = false

func init() {
	InitializeShell()
	EnableGlobalOptions()

	// Ensure some basic categories are created first so they
	// display consistently ahead of add-on command modules
	ensureCategory(CategoryHttp)
	ensureCategory(CategoryBenchmarks)
	ensureCategory(CategoryAnalysis)
	ensureCategory(CategoryUtilities)

	addCommands()
}

func addCommands() {
	AddCommand("rem", CategoryUtilities, NewRemCommand())
	AddCommand("run", CategoryUtilities, NewRunCommand())
	AddCommand("quit", CategoryUtilities, nil)
}

// InitializeShell -- Initialize common parameters needed by the shell
func InitializeShell() {
	if initialized {
		return
	}

	initialized = true
	curdir, err := os.Getwd()
	if err == nil {
		InitDirectory = curdir
		if len(curdir) > 0 && strings.HasSuffix(curdir, "/") == false {
			InitDirectory = InitDirectory + "/"
		}
	}

	exPath := ""
	{
		ex, err := os.Executable()
		if err == nil {
			exPath = filepath.Dir(ex)
		}
	}
	ExecutableDirectory = exPath
}

func GetInitDirectory() string {
	return InitDirectory
}

func GetExeDirectory() string {
	return ExecutableDirectory
}
