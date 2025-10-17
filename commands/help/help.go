package help

import (
	"fmt"
	"os"
	"sort"

	"example.com/commands"
	"example.com/commands/flag"
	"example.com/commands/registry"
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

	fmt.Fprintln(out)

	fmt.Fprintln(out, "Command flags:")
	for _, cmd := range commands.All() {
		specs := append([]flag.FlagSpec(nil), cmd.Flags...)
		sort.Slice(specs, func(i, j int) bool {
			return specs[i].Name < specs[j].Name
		})
		if len(specs) == 0 {
			fmt.Fprintf(out, "  %-10s (no flags)\n", cmd.Name)
			continue
		}
		fmt.Fprintf(out, "  %-10s\n", cmd.Name)
		for _, spec := range specs {
			fmt.Fprintf(out, "    --%-12s %s\n", spec.Name, spec.Usage)
		}
	}

	names := registry.SystemNames()
	if len(names) > 0 {
		sort.Strings(names)
		fmt.Fprintln(out, "\nSystems:")
		for _, name := range names {
			fmt.Fprintf(out, "  %-10s\n", name)
		}
	}

	return nil
}
