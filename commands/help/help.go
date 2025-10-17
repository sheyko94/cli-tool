package help

import (
	"fmt"
	"os"
	"sort"

	"example.com/commands"
	"example.com/internal/registry"
)

func init() {
	commands.Register(commands.Command{
		Name:        "help",
		Description: "Show usage information and available commands",
		Handler:     run,
	})
}

func run(ctx commands.Context) error {
	out := ctx.Stdout
	if out == nil {
		out = os.Stdout
	}

	if ctx.Usage != nil {
		ctx.Usage()
	}

	fmt.Fprintln(out, "\nAvailable commands:")
	for _, cmd := range commands.All() {
		fmt.Fprintf(out, "  %-10s %s\n", cmd.Name, cmd.Description)
	}

	names := registry.SystemNames()
	if len(names) > 0 {
		fmt.Fprintln(out, "\nSystems:")
		for _, name := range names {
			fmt.Fprintf(out, "  %-10s\n", name)
		}
	}

	specs := registry.FlagSpecs()
	if len(specs) > 0 {
		sort.Slice(specs, func(i, j int) bool {
			return specs[i].Name < specs[j].Name
		})
		fmt.Fprintln(out, "\nFlags:")
		for _, spec := range specs {
			fmt.Fprintf(out, "  -%-12s %s\n", spec.Name, spec.Usage)
		}
	}

	return nil
}
