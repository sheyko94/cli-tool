package github

import (
	"fmt"
	"strings"

	"example.com/internal/flag"
	"example.com/internal/systems/system"
)

const Name = "github"

type GitHubSystem struct{}

func (GitHubSystem) Name() string {
	return Name
}

func (GitHubSystem) Flags() []flag.FlagSpec {
	return []flag.FlagSpec{
		flag.NewStringFlag("github-repo", "", "GitHub repository to target (e.g., user/repo)"),
	}
}

func (GitHubSystem) Execute(ctx system.SystemContext) (string, error) {
	repoPtr, ok := ctx.Flags.String("github-repo")
	if !ok || strings.TrimSpace(*repoPtr) == "" {
		return "", fmt.Errorf("[GitHub] github-repo flag is required")
	}

	repo := strings.TrimSpace(*repoPtr)
	action := "summary"
	if len(ctx.Args) > 0 {
		action = strings.Join(ctx.Args, " ")
	}

	return fmt.Sprintf("[GitHub] Pretend to %s for repo %s...", action, repo), nil
}
