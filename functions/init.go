package functions

import (
	"github.com/brada954/restshell/shell"
)

func init() {
	// Register substitutes
	shell.RegisterSubstitutionHandler(NewGuidDefinition)
	shell.RegisterSubstitutionHandler(ToLowerDefinition)
	shell.RegisterSubstitutionHandler(ToUpperDefinition)
	shell.RegisterSubstitutionHandler(GetDateDefinition)
	shell.RegisterSubstitutionHandler(SetDateDefinition)

}
