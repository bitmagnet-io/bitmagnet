package cmd

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

type helpWriter struct {
	io.Writer
	err error
}

func (wr *helpWriter) printf(str string, args ...any) {
	if wr.err != nil {
		return
	}

	_, wr.err = fmt.Fprintf(wr, str, args...)
}

func (wr *helpWriter) print(str string) {
	if wr.err != nil {
		return
	}

	_, wr.err = fmt.Fprint(wr, str)
}

func (i *instance) printUsage(wr io.Writer) error {
	hw := &helpWriter{
		Writer: wr,
	}

	var cmdPath []string

	for curr := i; curr != nil; curr = i.parent {
		cmdPath = append(cmdPath, curr.Spec.Name)
	}

	slices.Reverse(cmdPath)

	hw.print("Usage: " + strings.Join(cmdPath, " "))

	for _, paramKey := range i.Spec.paramKeys {
		param := i.Spec.params[paramKey]
		shortLong := [2]string{param.Abbr, param.Name}
		for i, name := range shortLong {
			if name == "" {
				continue
			}
			tmpl := strings.Repeat("-", i+1) + name
			if param.Type != paramTypeBool {
				tmpl += " " + param.Placeholder
			}
			if param.Required {
				tmpl = "[" + tmpl + "]"
			}
			hw.print(" " + tmpl)
		}
	}

	// positionals?

	if len(i.Subcommands()) > 0 {
		hw.print(" <command> [<args>]")
	}

	return hw.err
}
