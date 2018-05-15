package shell

var ProgramName = "RestShell"
var RestShellInitFile = ".rsconfig"
var RestShellUserInitFile = ".rsconfig.user"

func init() {
	InitializeShell()
	EnableGlobalOptions()

	// Ensure some basic categories are created first so they
	// display consistently ahead of add-on command modules
	ensureCategory(CategoryHttp)
	ensureCategory(CategoryBenchmarks)
	ensureCategory(CategoryAnalysis)
	ensureCategory(CategoryUtilities)

	AddCommands()
}

func AddCommands() {
	AddCommand("rem", CategoryUtilities, NewRemCommand())
	AddCommand("run", CategoryUtilities, NewRunCommand())
	AddCommand("quit", CategoryUtilities, nil)
}
