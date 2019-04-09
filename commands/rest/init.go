package rest

import "github.com/brada954/restshell/shell"

func init() {
	AddCommands()
}

func AddCommands() {
	shell.AddCommand("base", shell.CategoryHttp, NewBaseCommand())
	shell.AddCommand("get", shell.CategoryHttp, NewGetCommand())
	shell.AddCommand("post", shell.CategoryHttp, NewPostCommand())
	shell.AddCommand("bmget", shell.CategoryBenchmarks, NewBmGetCommand())
	shell.AddCommand("bmpost", shell.CategoryBenchmarks, NewBmPostCommand())
	shell.AddCommand("smget", shell.CategoryBenchmarks, NewSmGetCommand())
	shell.AddCommand("login", shell.CategoryHttp, NewLoginCommand())
}
