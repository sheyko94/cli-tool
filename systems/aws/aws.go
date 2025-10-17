package aws

import (
	"fmt"
	"strings"

	"example.com/commands/flag"
	"example.com/systems/system"
)

const Name = "aws"

type AwsSystem struct{}

func (AwsSystem) Name() string {
	return Name
}

func (AwsSystem) Flags() []flag.FlagSpec {
	return []flag.FlagSpec{
		flag.NewStringFlag("service", "", "AWS service to inspect (e.g., rds, ec2)"),
	}
}

func (AwsSystem) Execute(ctx system.SystemContext) (string, error) {
	servicePtr, ok := ctx.Flags.String("service")
	if !ok || strings.TrimSpace(*servicePtr) == "" {
		return "", fmt.Errorf("[AWS] service flag is required")
	}

	service := strings.TrimSpace(*servicePtr)
	if len(ctx.Args) > 0 {
		return "", fmt.Errorf("[AWS] unexpected arguments: %s", strings.Join(ctx.Args, " "))
	}

	return fmt.Sprintf("[AWS] %s reports 3 healthy resources.", service), nil
}
