package registry

import (
	"errors"
	stdflag "flag"
	"fmt"
	"sort"
	"strings"
	"sync"

	"example.com/internal/flag"
	"example.com/internal/systems/aws"
	"example.com/internal/systems/github"
	"example.com/internal/systems/system"
)

var (
	mu          sync.RWMutex
	systems     = map[string]system.System{}
	systemNames []string
	fSpecs      []flag.FlagSpec
)

func init() {
	register(aws.AwsSystem{})
	register(github.GitHubSystem{})
}

// ErrSystemNotFound indicates that a system lookup failed.
var ErrSystemNotFound = errors.New("registry: system not registered")

// register makes a system discoverable and records its flags.
func register(system system.System) {
	mu.Lock()
	defer mu.Unlock()

	name := strings.TrimSpace(system.Name())
	if name == "" {
		panic("registry: system name is required")
	}
	if _, exists := systems[name]; exists {
		panic(fmt.Sprintf("registry: system %q already registered", name))
	}

	systems[name] = system
	systemNames = append(systemNames, name)
	registerFlagSpecsLocked(system.Flags())
}

// Get returns a registered system by name.
func Get(name string) (system.System, bool) {
	mu.RLock()
	defer mu.RUnlock()
	s, ok := systems[name]
	return s, ok
}

// SystemNames returns a sorted list of registered system names.
func SystemNames() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, len(systemNames))
	copy(names, systemNames)
	sort.Strings(names)
	return names
}

// Execute runs the named system with the provided context.
func Execute(name string, ctx system.SystemContext) (string, error) {
	mu.RLock()
	sys, ok := systems[name]
	mu.RUnlock()
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrSystemNotFound, name)
	}
	return sys.Execute(ctx)
}

// RegisterFlagProvider lets packages register flag specs that are not tied to a system.
func RegisterFlagProvider(provider flag.FlagProvider) {
	mu.Lock()
	defer mu.Unlock()
	registerFlagSpecsLocked(provider.Flags())
}

// FlagSpecs returns a copy of all registered flag specs.
func FlagSpecs() []flag.FlagSpec {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]flag.FlagSpec, len(fSpecs))
	copy(out, fSpecs)
	return out
}

// BindFlags binds known flag specs to the provided FlagSet and returns the stored pointers.
func BindFlags(fs *stdflag.FlagSet) flag.BoundFlags {
	mu.RLock()
	defer mu.RUnlock()

	b := flag.BoundFlags{
		StringFlags: make(map[string]*string, len(fSpecs)),
		BoolFlags:   make(map[string]*bool, len(fSpecs)),
	}

	for _, spec := range fSpecs {
		switch spec.Type {
		case flag.FlagString:
			def, _ := spec.Default.(string)
			ptr := fs.String(spec.Name, def, spec.Usage)
			b.StringFlags[spec.Name] = ptr
		case flag.FlagBool:
			def, _ := spec.Default.(bool)
			ptr := fs.Bool(spec.Name, def, spec.Usage)
			b.BoolFlags[spec.Name] = ptr
		default:
			// Future flag types can be added here.
		}
	}

	return b
}

func registerFlagSpecsLocked(specs []flag.FlagSpec) {
	for _, spec := range specs {
		if spec.Name == "" {
			panic("registry: flag name is required")
		}
		if hasFlagLocked(spec.Name) {
			panic(fmt.Sprintf("registry: flag %q already registered", spec.Name))
		}
		fSpecs = append(fSpecs, spec)
	}
}

func hasFlagLocked(name string) bool {
	for _, spec := range fSpecs {
		if spec.Name == name {
			return true
		}
	}
	return false
}
