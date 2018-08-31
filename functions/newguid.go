package functions

import (
	"github.com/brada954/restshell/shell"
	"github.com/satori/go.uuid"
)

func init() {
	shell.RegisterSubstitutionHandler(NewGuidDefinition)
}

var NewGuidDefinition = shell.SubstitutionFunction{
	Name:              "newguid",
	Group:             "newguid",
	FunctionHelp:      "Generate a guid formated as specified",
	Formats:           nil,
	OptionDescription: "",
	Options:           nil,
	Function:          NewGuidSubstitute,
}

// NewGuidSubstitute -- Implementatino of guid substitution
func NewGuidSubstitute(cache interface{}, subname string, format string, option string) (value string, data interface{}) {
	var guid uuid.UUID

	if cache == nil {
		var err error
		if guid, err = uuid.NewV4(); err != nil {
			guid = uuid.Nil
		}
	} else {
		guid = cache.(uuid.UUID)
	}

	switch format {
	default:
		return guid.String(), guid
	}
}
