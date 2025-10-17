package version

import (
	"fmt"
	"os"

	"example.com/commands"
)

func init() {
	commands.Register(commands.Command{
		Name:        "version",
		Description: "Print the application version",
		Handler:     run,
	})
}

func run(ctx commands.Context) error {
	out := ctx.Stdout
	if out == nil {
		out = os.Stdout
	}
	fmt.Fprintln(out, ctx.Version)
	return nil
}
