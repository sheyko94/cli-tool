package aws

import (
	"fmt"
	"os"

	"example.com/commands"
	"example.com/commands/registry"
	"example.com/systems/aws"
	"example.com/systems/system"
)

func init() {
	specs, ok := registry.SystemFlagSpecs(aws.Name)
	if !ok {
		panic("aws command: system not registered")
	}

	commands.Register(commands.Command{
		Name:        aws.Name,
		Description: "Interact with AWS resources",
		Flags:       specs,
		Handler:     run,
	})
}

func run(ctx commands.Context) error {
	out := ctx.Stdout
	if out == nil {
		out = os.Stdout
	}

	result, err := registry.Execute(aws.Name, system.SystemContext{
		Flags: ctx.Flags,
		Args:  ctx.Args,
	})
	if err != nil {
		return err
	}

	fmt.Fprintln(out, result)
	return nil
}
