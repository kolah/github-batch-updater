package github

import "context"

type Repository struct {
	Owner         string
	Name          string
	DefaultBranch Branch
}

type Branch struct {
	Name      string
	SHA       string
	IsDefault bool
}

type RepositoryService interface {
	GetRepository(ctx context.Context, owner, name string) (*Repository, error)
	CreateBranchFromDefaultBranch(ctx context.Context, repo *Repository, branchName string) error
}

func LoadRepository(owner, name, defaultBranchName string, defaultBranchSHA string) *Repository {
	return &Repository{
		Owner: owner,
		Name:  name,
		DefaultBranch: Branch{
			Name:      defaultBranchName,
			IsDefault: true,
			SHA:       defaultBranchSHA,
		},
	}
}
