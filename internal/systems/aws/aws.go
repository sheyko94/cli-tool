package aws

import (
	"fmt"
	"strings"

	"example.com/internal/flag"
	"example.com/internal/systems/system"
)

const Name = "aws"

type AwsSystem struct{}

func (AwsSystem) Name() string {
	return Name
}

func (AwsSystem) Flags() []flag.FlagSpec {
	return []flag.FlagSpec{
		flag.NewStringFlag("aws-service", "", "AWS service to target (e.g., rds, ec2)"),
	}
}

func (AwsSystem) Execute(ctx system.SystemContext) (string, error) {
	servicePtr, ok := ctx.Flags.String("aws-service")
	if !ok || strings.TrimSpace(*servicePtr) == "" {
		return "", fmt.Errorf("[AWS] aws-service flag is required")
	}

	service := strings.TrimSpace(*servicePtr)
	action := "describe"
	if len(ctx.Args) > 0 {
		action = strings.Join(ctx.Args, " ")
	}

	return fmt.Sprintf("[AWS] Pretend to %s for service %s...", action, service), nil
}
