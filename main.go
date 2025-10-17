package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"example.com/commands"
	_ "example.com/commands/aws"
	_ "example.com/commands/github"
	_ "example.com/commands/help"
	"example.com/commands/registry"
	_ "example.com/commands/version"
)

var version = "0.1.0"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr *os.File) int {
	if len(args) == 0 {
		args = []string{"help"}
	}

	cmdName := args[0]
	cmd, ok := commands.Lookup(cmdName)
	if !ok {
		fmt.Fprintf(stderr, "unknown command %q\n\n", cmdName)
		printGeneralUsage(stderr)
		return 1
	}

	fs := flag.NewFlagSet(cmd.Name, flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		printCommandUsage(stderr, cmd)
	}

	bound := registry.BindFlags(fs, cmd.Flags)
	if err := fs.Parse(args[1:]); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}

	ctx := commands.Context{
		Args:    fs.Args(),
		Flags:   bound,
		Usage:   func() { printGeneralUsage(stdout) },
		Stdout:  stdout,
		Stderr:  stderr,
		Version: version,
	}

	if err := cmd.Handler(ctx); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	return 0
}

func printGeneralUsage(out io.Writer) {
	if out == nil {
		out = os.Stdout
	}

	fmt.Fprintf(out, "Usage: %s <command> [flags] [args]\n", os.Args[0])
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Commands:")
	for _, cmd := range commands.All() {
		fmt.Fprintf(out, "  %-10s %s\n", cmd.Name, cmd.Description)
	}
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Use '<command> -h' to see command-specific flags.")
}

func printCommandUsage(out io.Writer, cmd commands.Command) {
	if out == nil {
		out = os.Stderr
	}

	fmt.Fprintf(out, "Usage: %s %s [flags] [args]\n", os.Args[0], cmd.Name)
	if len(cmd.Flags) == 0 {
		fmt.Fprintln(out)
		fmt.Fprintln(out, "This command does not accept any flags.")
		return
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, "Flags:")
	for _, spec := range cmd.Flags {
		fmt.Fprintf(out, "  --%-12s %s\n", spec.Name, spec.Usage)
	}
}
