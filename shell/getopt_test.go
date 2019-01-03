package shell

import "testing"

func TestCmdSetInterface(t *testing.T) {
	i := NewCmdSet()

	i.Reset()

	i.Args()
}

func TestStringList(t *testing.T) {
	i := NewCmdSet()

	list := i.StringListLong("header", 0, "A header name=value")

	CmdParse(i, []string{"cmd", "--header", "x=abc", "--header=y=123"})

	if len(list.Values) != 2 {
		t.Errorf("Invalid number of arguments: 2!=%d", len(list.Values))
	}

	f := func(t *testing.T, expected string, value string) {
		if expected != value {
			t.Errorf("First arg was not equal: x=abc!=%s", list.Values[0])
		}
	}

	if len(list.Values) > 0 {
		f(t, "x=abc", list.Values[0])
	}

	if len(list.Values) > 1 {
		f(t, "y=123", list.Values[1])
	}
}
