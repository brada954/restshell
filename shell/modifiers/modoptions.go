package modifiers

import (
	"github.com/brada954/restshell/shell"
)

// ModifierOptions -- Common options for modifiers
type ModifierOptions struct {
	toLowerOption    *bool
	toUpperOption    *bool
	regexOption      *string
	strToIntOption   *bool
	strToFloatOption *bool
	lenOption        *bool
	quoteOption      *bool
}

// AddModifierOptions -- Add options for modifiers
func AddModifierOptions(set shell.CmdSet) ModifierOptions {
	options := ModifierOptions{}
	options.toLowerOption = set.BoolLong("to-lower", 0, "Convert string value to lowercase")
	options.toUpperOption = set.BoolLong("to-upper", 0, "Convert string value to uppercase")
	options.regexOption = set.StringLong("regex", 0, "", "Use regex to extract from value (1st)")
	options.strToIntOption = set.BoolLong("int", 0, "Convert string value to int (2nd)")
	options.strToFloatOption = set.BoolLong("float", 0, "Convert string value to float (3rd)")
	options.lenOption = set.BoolLong("len", 0, "Use the length of the value (4th)")
	options.quoteOption = set.BoolLong("quote", 0, "Quote a string and escape internal quotes")
	return options
}

// ConstructModifier -- Build a modifier chain based on common options
func ConstructModifier(options ModifierOptions) ValueModifier {

	// Order is a set precendence; changing order can be catastrophic. Generally, the
	// modifiers call the previous modifier before performing its own task.
	//
	// For example, conversion of string to int is initialized before string to float, so
	// a string containing float can be converted to float and then to int. A "float string",
	// will fail to convert to int because of the period in the text string.

	valueModifierFunc := NullModifier
	if *options.toLowerOption {
		valueModifierFunc = MakeStringToLowerModifier(valueModifierFunc)
	} else if *options.toUpperOption {
		valueModifierFunc = MakeStringToUpperModifier(valueModifierFunc)
	}
	if len(*options.regexOption) > 0 {
		valueModifierFunc = MakeRegExModifier(*options.regexOption, valueModifierFunc)
	}
	if *options.quoteOption {
		valueModifierFunc = MakeQuoteModifier(valueModifierFunc)
	}
	if *options.strToFloatOption {
		valueModifierFunc = MakeToFloatModifier(valueModifierFunc)
	}
	if *options.strToIntOption {
		valueModifierFunc = MakeToIntModifier(valueModifierFunc)
	}
	if *options.lenOption {
		valueModifierFunc = MakeLengthModifier(valueModifierFunc)
	}

	return valueModifierFunc
}
