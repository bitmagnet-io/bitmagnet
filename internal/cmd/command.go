package cmd

import (
	"errors"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/env"
)

type Factory[C Command] func() C

type CommandFactory = Factory[Command]

func (f Factory[C]) Run(env env.Env) (int, error) {
	cmd := f()

	spec, err := introspect(cmd)
	if err != nil {
		return 1, fmt.Errorf("%w: %w", Err, err)
	}

	inst := spec.newInstance(cmd, env.Args())

	err = inst.run(env)

	if err != nil {
		exitCode, err := inst.OnError(env, err)

		if err != nil {
			err = fmt.Errorf("%w: %w", Err, err)
		}

		return exitCode, err
	}

	return 0, nil
}

type Command interface {
	Setup(env.Env) error
	Subcommands() []Command
	Run(env.Env) error
	Teardown(env.Env) error
	OnError(env.Env, error) (int, error)
}

type Cmd struct {
	*instance
}

func (*Cmd) Setup(env.Env) error {
	return nil
}

func (*Cmd) Subcommands() []Command {
	return nil
}

func (c *Cmd) Run(env env.Env) error {
	return c.printUsage(env)
}

func (*Cmd) Teardown(env.Env) error {
	return nil
}

func (*Cmd) OnError(env env.Env, err error) (int, error) {
	fmt.Fprintln(env.Stderr(), err)
	switch {
	case errors.Is(err, ErrInvalidArgs):
		return 2, err
	default:
		return 1, err
	}
}

var _ Command = (*Cmd)(nil)
