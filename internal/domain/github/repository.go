package github

import "context"

type Repository struct {
	Owner        string
	Name         string
	SourceBranch Branch
	TargetBranch Branch
}

type Branch struct {
	Name string
	SHA  string
}

func (b Branch) Empty() bool {
	return b.Name == "" && b.SHA == ""
}

func (b Branch) Equal(other Branch) bool {
	return b.Name == other.Name && b.SHA == other.SHA
}

type RepositoryService interface {
	GetRepository(ctx context.Context, owner, name, sourceBranch string) (*Repository, error)
	CreateBranchFromSourceBranch(ctx context.Context, repo *Repository, branchName string) error
}

func LoadRepository(owner, name, branchName string, branchSHA string) *Repository {
	return &Repository{
		Owner: owner,
		Name:  name,
		SourceBranch: Branch{
			Name: branchName,
			SHA:  branchSHA,
		},
	}
}
