package system

import "example.com/internal/flag"

// SystemContext represents the runtime context passed to a system execution.
type SystemContext struct {
	Flags flag.BoundFlags
	Args  []string
}

// System defines the contract each system must satisfy.
type System interface {
	Name() string
	Flags() []flag.FlagSpec
	Execute(SystemContext) (string, error)
}
