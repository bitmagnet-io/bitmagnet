package spec

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"google.golang.org/protobuf/compiler/protogen"
)

func fileVarName(f *protogen.File, suffix string) string {
	prefix := f.GoDescriptorIdent.GoName
	_, n := utf8.DecodeRuneInString(prefix)
	prefix = strings.ToLower(prefix[:n]) + prefix[n:]
	return prefix + "_" + suffix
}

// parseParam parses a comment and extract parameters for go-plugin
// e.g. // go:plugin type=plugin version=2
func parseParam(comment string) (Parameter, error) {
	param := Parameter{
		APIVersion: 1,
		Type:       ServiceNone,
		Module:     EnvModuleName,
	}
	for _, line := range strings.Split(comment, "\n") {
		line = strings.TrimPrefix(line, "//")
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "go:plugin") {
			continue
		}
		line = strings.TrimPrefix(line, "go:plugin")

		for _, field := range strings.Fields(line) {
			var key, value string
			if i := strings.Index(field, "="); i >= 0 {
				value = field[i+1:]
				key = field[0:i]
			}
			if key == "" || value == "" {
				continue
			}
			switch key {
			case "type":
				switch value {
				case "host":
					param.Type = ServiceHost
				case "plugin":
					param.Type = ServicePlugin
				default:
					param.Type = ServiceUnknown
				}
			case "version":
				ver, err := strconv.Atoi(value)
				if err != nil {
					return Parameter{}, fmt.Errorf("invalid version: %w", err)
				}
				param.APIVersion = ver
			case "module":
				if param.Type == ServiceHost && len(value) > 0 {
					param.Module = value
				}
			}
		}
	}
	return param, nil
}
