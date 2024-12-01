package batch

import (
	"fmt"

	"github.com/kolah/github-batch-updater/internal/application"
	githubDomain "github.com/kolah/github-batch-updater/internal/domain/github"
	transformerDomain "github.com/kolah/github-batch-updater/internal/domain/transformer"
	"github.com/kolah/github-batch-updater/internal/infrastructure/transformer"
	"github.com/kolah/github-batch-updater/internal/pkg/slices"
	"gopkg.in/yaml.v3"
)

type Repository struct {
	Owner string
	Name  string
}

type Replace struct {
	Match   string
	Replace string
}

type Replaces struct {
	Replaces []Replace
}

func (f Replaces) Type() string {
	return "replace"
}

type YQEdit struct {
	Expressions []string
}

func (f YQEdit) Type() string {
	return "yq"
}

type PullRequest struct {
	Title         string
	Body          string
	Reviewers     []string
	TeamReviewers []string `yaml:"team_reviewers"`
}

type Operation interface {
	Type() string
}

type FileOperation struct {
	Operation Operation
}

type Config struct {
	Repositories []Repository
	Files        map[string][]FileOperation
	CreatePR     bool         `yaml:"create_pr"`
	PullRequest  *PullRequest `yaml:"pull_request"`
	TargetBranch string       `yaml:"target_branch"`
}

func (fo *FileOperation) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("expected mapping node for operation")
	}

	if len(node.Content) != 2 {
		return fmt.Errorf("expected operation to have 2 nodes")
	}

	var kind string

	if err := node.Content[0].Decode(&kind); err != nil {
		return err
	}
	switch kind {
	case "replace":
		var fr Replaces
		if err := node.Content[1].Decode(&fr); err != nil {
			return err
		}

		fo.Operation = fr
		return nil
	case "yq":
		var yq YQEdit
		if err := node.Content[1].Decode(&yq); err != nil {
			return err
		}

		fo.Operation = yq
		return nil
	}

	return fmt.Errorf("unknown operation type: %s", kind)
}

func (c *Config) ToInput(dryRun bool) application.Input {
	rules := make(map[string][]transformerDomain.Rule)
	for file, rulesConfig := range c.Files {
		for _, rule := range rulesConfig {
			switch r := rule.Operation.(type) {
			case Replaces:
				rules[file] = append(rules[file], transformer.ReplaceRule{
					Replaces: slices.Map(r.Replaces, func(r Replace) transformer.Replace {
						return transformer.Replace{
							Match:   r.Match,
							Replace: r.Replace,
						}
					}),
				})
			case YQEdit:
				rules[file] = append(rules[file], transformer.YQRule{
					Expressions: r.Expressions,
				})
			}

		}
	}
	op := application.Input{
		Repositories: slices.Map(c.Repositories, func(r Repository) application.RepositoryInput {
			return application.RepositoryInput{
				Owner: r.Owner,
				Name:  r.Name,
			}
		}),
		Rules: rules,
		PullRequestTemplate: githubDomain.PullRequestTemplate{
			Title:         c.PullRequest.Title,
			Body:          c.PullRequest.Body,
			Reviewers:     c.PullRequest.Reviewers,
			TeamReviewers: c.PullRequest.TeamReviewers,
		},
		TargetBranch: c.TargetBranch,
		CreatePR:     c.CreatePR,
		DryRun:       dryRun,
	}

	return op
}
