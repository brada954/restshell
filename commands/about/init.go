package about

import "github.com/brada954/restshell/shell"

func init() {
	addCommands()
	addAboutTopics()
}

func addCommands() {
	shell.AddCommand("about", shell.CategoryHelp, NewAboutCommand())
	shell.AddCommand("version", shell.CategoryHelp, NewVersionCommand())
	shell.AddCommand("help", shell.CategoryHelp, nil)
}

func addAboutTopics() {
	shell.AddAboutTopic(NewAuthTopic())
	shell.AddAboutTopic(NewBenchmarkTopic())
	shell.AddAboutTopic(NewJsonPathTopic())
	shell.AddAboutTopic(NewSubstitutionTopic())
}
