package cmd

import (
	"bytes"
	"testing"

	"github.com/bitmagnet-io/bitmagnet/internal/env"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCommand struct {
	Cmd             `cmd:"name=testing"`
	EnabledPlugins  CSV[ref.Ref]
	DisabledPlugins CSVStringSlice
	sub             *string
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
			"--enabled-plugins=foo.bar,baz.bat",
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
		cmd.EnabledPlugins,
	)

	assert.Equal(t, "sub", str)

	var buff bytes.Buffer

	inst.printUsage(&buff)

	println(buff.String())
}
