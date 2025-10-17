package flag

// FlagType represents the kind of command-line flag
type FlagType int

const (
	FlagString FlagType = iota
	FlagBool
)

// FlagSpec describes a single flag to expose through the CLI.
type FlagSpec struct {
	Name    string
	Usage   string
	Type    FlagType
	Default any
}

// BoundFlags contains the results of binding FlagSpecs onto a FlagSet.
type BoundFlags struct {
	StringFlags map[string]*string
	BoolFlags   map[string]*bool
}

// NewStringFlag creates a string flag spec.
func NewStringFlag(name, defaultValue, usage string) FlagSpec {
	return FlagSpec{
		Name:    name,
		Type:    FlagString,
		Default: defaultValue,
		Usage:   usage,
	}
}

// NewBoolFlag creates a bool flag spec.
func NewBoolFlag(name string, defaultValue bool, usage string) FlagSpec {
	return FlagSpec{
		Name:    name,
		Type:    FlagBool,
		Default: defaultValue,
		Usage:   usage,
	}
}

// String retrieves a registered string flag pointer by name.
func (b BoundFlags) String(name string) (*string, bool) {
	ptr, ok := b.StringFlags[name]
	return ptr, ok
}

// Bool retrieves a registered bool flag pointer by name.
func (b BoundFlags) Bool(name string) (*bool, bool) {
	ptr, ok := b.BoolFlags[name]
	return ptr, ok
}

// FlagProvider allows packages to register global (non-system) flags.
type FlagProvider interface {
	Flags() []FlagSpec
}
