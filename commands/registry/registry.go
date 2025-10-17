package registry

import (
	"errors"
	stdflag "flag"
	"fmt"
	"sort"
	"strings"
	"sync"

	"example.com/commands/flag"
	"example.com/systems/aws"
	"example.com/systems/github"
	"example.com/systems/system"
)

var (
	mu       sync.RWMutex
	systems  = map[string]registeredSystem{}
	ordering []string
)

type registeredSystem struct {
	impl  system.System
	flags []flag.FlagSpec
}

func init() {
	MustRegisterSystem(aws.AwsSystem{})
	MustRegisterSystem(github.GitHubSystem{})
}

var ErrSystemNotFound = errors.New("registry: system not registered")

func MustRegisterSystem(sys system.System) {
	if err := RegisterSystem(sys); err != nil {
		panic(err)
	}
}

func RegisterSystem(sys system.System) error {
	name := strings.TrimSpace(sys.Name())
	if name == "" {
		return fmt.Errorf("registry: system name is required")
	}

	mu.Lock()
	defer mu.Unlock()

	key := strings.ToLower(name)
	if _, exists := systems[key]; exists {
		return fmt.Errorf("registry: system %q already registered", name)
	}

	systems[key] = registeredSystem{
		impl:  sys,
		flags: cloneFlagSpecs(sys.Flags()),
	}
	ordering = append(ordering, key)
	sort.Strings(ordering)

	return nil
}

func Execute(name string, ctx system.SystemContext) (string, error) {
	sys, err := lookup(name)
	if err != nil {
		return "", err
	}
	return sys.impl.Execute(ctx)
}

func SystemNames() []string {
	mu.RLock()
	defer mu.RUnlock()

	names := make([]string, len(ordering))
	copy(names, ordering)
	return names
}

func SystemFlagSpecs(name string) ([]flag.FlagSpec, bool) {
	sys, err := lookup(name)
	if err != nil {
		return nil, false
	}
	return cloneFlagSpecs(sys.flags), true
}

func BindFlags(fs *stdflag.FlagSet, specs []flag.FlagSpec) flag.BoundFlags {
	bound := flag.BoundFlags{
		StringFlags: make(map[string]*string, len(specs)),
		BoolFlags:   make(map[string]*bool, len(specs)),
	}

	for _, spec := range specs {
		switch spec.Type {
		case flag.FlagString:
			def, _ := spec.Default.(string)
			bound.StringFlags[spec.Name] = fs.String(spec.Name, def, spec.Usage)
		case flag.FlagBool:
			def, _ := spec.Default.(bool)
			bound.BoolFlags[spec.Name] = fs.Bool(spec.Name, def, spec.Usage)
		}
	}

	return bound
}

func lookup(name string) (registeredSystem, error) {
	key := strings.ToLower(strings.TrimSpace(name))
	if key == "" {
		return registeredSystem{}, ErrSystemNotFound
	}

	mu.RLock()
	sys, ok := systems[key]
	mu.RUnlock()
	if !ok {
		return registeredSystem{}, fmt.Errorf("%w: %s", ErrSystemNotFound, name)
	}
	return sys, nil
}

func cloneFlagSpecs(specs []flag.FlagSpec) []flag.FlagSpec {
	if len(specs) == 0 {
		return nil
	}
	cp := make([]flag.FlagSpec, len(specs))
	copy(cp, specs)
	return cp
}
