package cmd

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"

	"github.com/bitmagnet-io/bitmagnet/internal/env"
)

type instance struct {
	Spec
	Command
	parent        *instance
	values        map[string][]any
	cmdValue      reflect.Value
	reflectValues map[string]reflect.Value
	args          []string
	index         int
	subcmds       []Command
}

func (i *instance) Subcommands() []Command {
	if i.subcmds != nil {
		return i.subcmds
	}

	i.subcmds = i.Command.Subcommands()

	if i.subcmds == nil {
		i.subcmds = make([]Command, 0)
	}

	return i.subcmds
}

func (i *instance) run(env env.Env) error {
	cmd := i.Command

outer:
	for {
		token, ok := i.nextToken()
		if !ok {
			break
		}

		switch token.tokenType {
		case tokenTypeKeyValue:
			p, err := i.Spec.param(token.key)
			if err != nil {
				return fmt.Errorf("%w: %s: %w", ErrInvalidArgs, i.Name, err)
			}

			if err := i.handleParamValue(p, token.value); err != nil {
				return fmt.Errorf("%w: %s: %w", ErrInvalidArgs, i.Name, err)
			}

		case tokenTypeKey:
			p, err := i.Spec.param(token.key)
			if err != nil {
				return fmt.Errorf("%w: %s: %w", ErrInvalidArgs, i.Name, err)
			}

			if len(i.values[p.Name]) > 0 && !p.Multiple {
				return fmt.Errorf("%w: %s: %w: %s", ErrInvalidArgs, i.Name, ErrRepeatedParam, p.Name)
			}

			var (
				arg string
				ok  bool
			)

			if p.Type != paramTypeBool {
				arg, ok = i.nextArg()
				if !ok {
					return fmt.Errorf("%w: %s: %w: %s", ErrInvalidArgs, i.Name, ErrMissingParamValue, p.Name)
				}
			}

			if err := i.handleParamValue(p, arg); err != nil {
				return fmt.Errorf("%w: %s: %w", ErrInvalidArgs, i.Name, err)
			}
		case tokenTypePositional:
			i.index--

			break outer
		default:
			panic("unexpected token")
		}
	}

	i.applyValues()

	if err := cmd.Setup(env); err != nil {
		return fmt.Errorf("%w: %s: %w: %w", ErrExecution, i.Name, ErrSetup, err)
	}

	if token, ok := i.nextToken(); ok {
		if token.tokenType == tokenTypePositional {
			for _, subCmd := range i.Subcommands() {
				subCmdSpec, err := introspect(subCmd)
				if err != nil {
					return err
				}
				if subCmdSpec.Name != token.value {
					continue
				}
				subInstance := subCmdSpec.newInstance(subCmd, i.args[i.index:])
				subInstance.parent = i
				return subInstance.run(env)
			}
		}

		return fmt.Errorf("%w: %s: %w: %s", ErrInvalidArgs, i.Name, ErrUnexpectedArgument, token.arg)
	}

	err := cmd.Run(env)
	if err != nil {
		return fmt.Errorf("%w: %s: %w", ErrExecution, i.Name, err)
	}

	err = cmd.Teardown(env)
	if err != nil {
		return fmt.Errorf("%w: %s: %w: %w", ErrExecution, i.Name, ErrTeardown, err)
	}

	return nil
}

func (i *instance) nextToken() (token, bool) {
	arg, ok := i.nextArg()
	if !ok {
		return token{}, false
	}

	return parseToken(arg), true
}

func (i *instance) nextArg() (string, bool) {
	if i.index >= len(i.args) {
		return "", false
	}

	arg := i.args[i.index]

	i.index++

	return arg, true
}

func (i *instance) handleParamValue(param Param, strValue string) error {
	var (
		values []any
	)

	switch param.Type {
	case paramTypeTextUnmarshaler:
		v := i.reflectValues[param.Name]
		ptr := reflect.New(v.Type())
		ptr.Elem().Set(v)
		unmarshaler := ptr.Interface().(encoding.TextUnmarshaler)
		err := unmarshaler.UnmarshalText([]byte(strValue))
		if err != nil {
			return fmt.Errorf("%w: %s: %w", ErrParseParamValue, param.Name, err)
		}
		values = append(values, unmarshaler)
	case paramTypeBool:
		if strValue == "" {
			values = append(values, true)
		} else {
			value, err := strconv.ParseBool(strValue)
			if err != nil {
				return fmt.Errorf("%w: %s: %w", ErrParseParamValue, param.Name, err)
			}
			values = append(values, value)
		}
	case paramTypeString:
		values = append(values, strValue)
	}

	i.values[param.Name] = append(i.values[param.Name], values...)

	return nil
}

func (i *instance) applyValues() {
	for name, values := range i.values {
		param := i.Spec.params[name]
		var value any

		switch param.Type {
		case paramTypeBool:
			value = values[0].(bool)
		case paramTypeString:
			value = values[0].(string)
		case paramTypeTextUnmarshaler:
			value = reflect.ValueOf(values[0]).Elem().Interface()
		default:
			panic("unknown param type")
		}
		i.reflectValues[name].Set(reflect.ValueOf(value))
	}

	i.cmdValue.Set(reflect.ValueOf(Cmd{
		instance: i,
	}))
}
