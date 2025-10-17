package main

import (
	"flag"
	"fmt"

	"example.com/internal/registry"
)

func main() {
	boundFlags := registry.BindFlags(flag.CommandLine)

	fmt.Printf("flags received:\n")
	for name, value := range boundFlags.StringFlags {
		fmt.Printf("  %s: %v\n", name, *value)
	}

}
