package github

import (
	"fmt"
	"os"

	"example.com/commands"
	"example.com/commands/registry"
	"example.com/systems/github"
	"example.com/systems/system"
)

func init() {
	specs, ok := registry.SystemFlagSpecs(github.Name)
	if !ok {
		panic("github command: system not registered")
	}

	commands.Register(commands.Command{
		Name:        github.Name,
		Description: "Interact with GitHub repositories",
		Flags:       specs,
		Handler:     run,
	})
}

func run(ctx commands.Context) error {
	out := ctx.Stdout
	if out == nil {
		out = os.Stdout
	}

	result, err := registry.Execute(github.Name, system.SystemContext{
		Flags: ctx.Flags,
		Args:  ctx.Args,
	})
	if err != nil {
		return err
	}

	fmt.Fprintln(out, result)
	return nil
}
