package cmd

import (
	"bytes"
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/pkg/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCommand struct {
	Cmd     `cmd:"name=testing"`
	Refs    CSV[ref.Ref] `cmd:"doc=List of refs,example='foo.bar,baz.bat'"`
	Strings CSVStringSlice
	sub     *string
}

func (c *testCommand) Run(env.Env) error {
	return nil
}

func (c *testCommand) Subcommands() []Command {
	return []Command{&testSubCommand{
		sub: c.sub,
	}}
}

type testSubCommand struct {
	Cmd     `cmd:"name=sub"`
	TestArg string
	sub     *string
}

func (c *testSubCommand) Run(env.Env) error {
	*c.sub = c.TestArg

	return nil
}

func TestArg(t *testing.T) {
	t.Parallel()

	str := ""

	sub := &str

	cmd := &testCommand{
		sub: sub,
	}
	spec, err := introspect(cmd)
	require.NoError(t, err)

	env := env.New(
		env.WithContext(t.Context()),
		env.WithArgs(
			"--refs=foo.bar,baz.bat",
			"sub",
			"--test-arg=sub",
		),
	)

	inst := spec.newInstance(cmd, env.Args())
	err = inst.run(env)
	require.NoError(t, err)

	assert.Equal(
		t,
		CSV[ref.Ref]{
			ref.MustParse("foo.bar"),
			ref.MustParse("baz.bat"),
		},
		cmd.Refs,
	)

	assert.Equal(t, "sub", str)
}

func TestHelp(t *testing.T) {
	t.Parallel()

	str := ""

	sub := &str

	cmd := &testCommand{
		sub: sub,
	}
	spec, err := introspect(cmd)
	require.NoError(t, err)

	stdout := &bytes.Buffer{}

	env := env.New(
		env.WithContext(t.Context()),
		env.WithArgs("--help"),
		env.WithStdout(stdout),
	)

	inst := spec.newInstance(cmd, env.Args())
	err = inst.run(env)
	require.NoError(t, err)

	assert.Equal(t, `Usage:

testing --refs=foo.bar,baz.bat --strings=STRINGS <command> [<args>]

Parameters:

  --refs=foo.bar,baz.bat
                         List of refs
  --strings=STRINGS
  --help, -h             Show help for this command and exit

Commands:

  sub
`, stdout.String())
}
