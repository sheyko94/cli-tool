package github

import (
	"fmt"
	"strings"

	"example.com/commands/flag"
	"example.com/systems/system"
)

const Name = "github"

type GitHubSystem struct{}

func (GitHubSystem) Name() string {
	return Name
}

func (GitHubSystem) Flags() []flag.FlagSpec {
	return []flag.FlagSpec{
		flag.NewStringFlag("repository", "", "GitHub repository to inspect (e.g., user/repo)"),
	}
}

func (GitHubSystem) Execute(ctx system.SystemContext) (string, error) {
	repoPtr, ok := ctx.Flags.String("repository")
	if !ok || strings.TrimSpace(*repoPtr) == "" {
		return "", fmt.Errorf("[GitHub] repository flag is required")
	}

	if len(ctx.Args) > 0 {
		return "", fmt.Errorf("[GitHub] unexpected arguments: %s", strings.Join(ctx.Args, " "))
	}

	repo := strings.TrimSpace(*repoPtr)

	return fmt.Sprintf("[GitHub] %s has 15 open MRs.", repo), nil
}
