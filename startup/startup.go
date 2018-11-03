package startup

import (
	// Import the basic functionality of restshell; each package
	// has an init function to inject itself into the shell
	_ "github.com/brada954/restshell/commands/about"
	_ "github.com/brada954/restshell/commands/rest"
	_ "github.com/brada954/restshell/commands/result"
	_ "github.com/brada954/restshell/commands/util"
	_ "github.com/brada954/restshell/functions"
)
