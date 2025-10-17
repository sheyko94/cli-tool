package commands

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"

	"example.com/commands/flag"
)

// Context carries runtime information provided to each command execution.
type Context struct {
	Args    []string
	Flags   flag.BoundFlags
	Usage   func()
	Stdout  io.Writer
	Stderr  io.Writer
	Version string
}

// Handler is the function signature for command implementations.
type Handler func(Context) error

// Command describes a runnable CLI command.
type Command struct {
	Name        string
	Description string
	Flags       []flag.FlagSpec
	Handler     Handler
}

var (
	mu    sync.RWMutex
	reg   = map[string]Command{}
	order []string
)

// Register adds a new command for use in the CLI.
func Register(cmd Command) {
	if cmd.Name == "" {
		panic("commands: name is required")
	}
	if cmd.Handler == nil {
		panic(fmt.Sprintf("commands: handler is required for %q", cmd.Name))
	}

	name := strings.ToLower(cmd.Name)
	cmd = cloneCommand(cmd)

	mu.Lock()
	defer mu.Unlock()

	if _, exists := reg[name]; exists {
		panic(fmt.Sprintf("commands: command %q already registered", cmd.Name))
	}
	reg[name] = cmd
	order = append(order, cmd.Name)
	sort.Strings(order)
}

// Lookup retrieves a command by name (case-insensitive).
func Lookup(name string) (Command, bool) {
	mu.RLock()
	defer mu.RUnlock()
	cmd, ok := reg[strings.ToLower(name)]
	if !ok {
		return Command{}, false
	}
	return cloneCommand(cmd), true
}

// All returns a snapshot of all registered commands, sorted by name.
func All() []Command {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]Command, 0, len(order))
	for _, name := range order {
		cmd, ok := reg[strings.ToLower(name)]
		if ok {
			out = append(out, cloneCommand(cmd))
		}
	}
	return out
}

func cloneCommand(cmd Command) Command {
	if len(cmd.Flags) == 0 {
		return cmd
	}
	flags := make([]flag.FlagSpec, len(cmd.Flags))
	copy(flags, cmd.Flags)
	cmd.Flags = flags
	return cmd
}
