package cmd

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

const helpParam = "help"

type helpWriter struct {
	io.Writer
	err error
}

func (wr *helpWriter) Write(p []byte) (n int, err error) {
	if wr.err != nil {
		return 0, wr.err
	}
	n, err = wr.Writer.Write(p)
	if err != nil {
		wr.err = err
	}
	return n, err
}

func (wr *helpWriter) print(str string) {
	if wr.err != nil {
		return
	}

	_, wr.err = fmt.Fprint(wr.Writer, str)
}

// the width of the left column
const colWidth = 25

// printOption prints a line like this:
//
//	--option FOO            A description of the option [default: 123]
//
// If the text on the left is longer than a certain threshold, the description is moved to the next line:
//
//	--verylongoptionoption VERY_LONG_VARIABLE
//	                        A description of the option [default: 123]
//
// If multiple "extras" are provided then they are put inside a single set of square brackets:
//
//	--option FOO            A description of the option [default: 123, env: FOO]
func (wr *helpWriter) printOption(item, description string, bracketed ...string) {
	lhs := "  " + item
	wr.print(lhs)
	if description != "" {
		if len(lhs)+2 < colWidth {
			wr.print(strings.Repeat(" ", colWidth-len(lhs)))
		} else {
			wr.print("\n" + strings.Repeat(" ", colWidth))
		}
		wr.print(description)
	}

	var brack string
	for _, s := range bracketed {
		if s != "" {
			if brack != "" {
				brack += ", "
			}
			brack += s
		}
	}

	if brack != "" {
		wr.print(fmt.Sprintf(" [%s]", brack))
	}

	wr.print("\n")
}

func (wr *helpWriter) wrapError() error {
	if wr.err != nil {
		return fmt.Errorf("%w: %w", ErrHelp, wr.err)
	}

	return nil
}

func (i *instance) printHelp(wr io.Writer) error {
	hw := &helpWriter{
		Writer: wr,
	}

	var path []string
	for curr := i; curr != nil; curr = curr.parent {
		path = append(path, curr.Spec.name)
	}
	slices.Reverse(path)

	if len(path) > 1 || i.Spec.doc != "" {
		hw.print(strings.Join(path, " "))
		if i.Spec.doc != "" {
			hw.print(": " + i.Spec.doc)
		}

		hw.print("\n\n")
	}

	hw.print("Usage:\n\n" + strings.Join(path, " [<args>] "))
	for _, paramKey := range i.Spec.paramKeys {
		if paramKey == helpParam {
			continue
		}
		param := i.Spec.params[paramKey]
		shortLong := [2]string{param.Abbr, param.Name}
		for i, name := range shortLong {
			if name == "" {
				continue
			}
			tmpl := strings.Repeat("-", i+1) + name
			if param.Type != paramTypeBool {
				tmpl += "=" + param.Example
			}
			if param.Required {
				tmpl = "[" + tmpl + "]"
			}
			hw.print(" " + tmpl)
		}
	}

	subcmds := i.Subcommands()

	if len(subcmds) > 0 {
		hw.print(" <command> [<args>]")
	}

	hw.print("\n\n")

	if len(i.Spec.paramKeys) > 0 {
		hw.print("Parameters:\n\n")
		for _, paramKey := range i.Spec.paramKeys {
			param := i.Spec.params[paramKey]
			ways := make([]string, 0, 2)
			if param.Name != "" {
				ways = append(ways, synopsis(param, "--"+param.Name))
			}
			if param.Abbr != "" {
				ways = append(ways, synopsis(param, "-"+param.Abbr))
			}
			if len(ways) > 0 {
				hw.printOption(strings.Join(ways, ", "), param.Doc, withDefault(param.Default))
			}
		}
		hw.print("\n")
	}

	if len(subcmds) > 0 {
		hw.print("Commands:\n\n")
		for _, subcmd := range subcmds {
			spec, err := introspect(subcmd)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrHelp, err)
			}

			hw.print("  " + spec.name)
			if spec.doc != "" {
				hw.print(strings.Repeat(" ", colWidth-len(spec.name)-2) + spec.doc)
			}
			hw.print("\n")
		}
	}

	return hw.wrapError()
}

func synopsis(param Param, form string) string {
	// if the user omits the placeholder tag then we pick one automatically,
	// but if the user explicitly specifies an empty placeholder then we
	// leave out the placeholder in the help message
	if param.Type == paramTypeBool || param.Example == "" {
		return form
	}
	return form + "=" + param.Example
}

func withDefault(s string) string {
	if s == "" {
		return ""
	}
	return "default: " + s
}
