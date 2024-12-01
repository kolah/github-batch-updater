package application

import (
	"github.com/kolah/github-batch-updater/internal/domain/github"
	"github.com/kolah/github-batch-updater/internal/domain/transformer"
)

type Input struct {
	Repositories        []RepositoryInput
	Rules               map[string][]transformer.Rule
	TargetBranch        string
	PullRequestTemplate github.PullRequestTemplate
	CreatePR            bool
	DryRun              bool
}

type RepositoryInput struct {
	Owner string
	Name  string
}
