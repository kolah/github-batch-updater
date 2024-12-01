package infrastructure

import (
	"context"
	"net/http"

	"github.com/google/go-github/v60/github"
	githubDomain "github.com/kolah/github-batch-updater/internal/domain/github"
	"github.com/kolah/github-batch-updater/internal/pkg/errors"
)

type GithubRepositoryService struct {
	client *github.Client
}

func NewGithubRepositoryService(client *github.Client) *GithubRepositoryService {
	return &GithubRepositoryService{client: client}
}

func (g GithubRepositoryService) GetRepository(ctx context.Context, owner, name string) (*githubDomain.Repository, error) {
	repositoryResult, _, err := g.client.Repositories.Get(ctx, owner, name)

	if err != nil {
		return nil, g.handleRepositoryError(err)
	}

	ref, _, err := g.client.Git.GetRef(ctx, owner, name, "refs/heads/"+repositoryResult.GetDefaultBranch())
	if err != nil {
		return nil, g.handleRefError(err)
	}

	defaultBranchSHA := *ref.Object.SHA

	return githubDomain.LoadRepository(owner, name, repositoryResult.GetDefaultBranch(), defaultBranchSHA), nil
}

func (g GithubRepositoryService) CreateBranchFromDefaultBranch(
	ctx context.Context,
	repo *githubDomain.Repository,
	branchName string,
) error {
	_, _, err := g.client.Git.CreateRef(ctx, repo.Owner, repo.Name, &github.Reference{
		Ref: github.String("refs/heads/" + branchName),
		Object: &github.GitObject{
			SHA: github.String(repo.DefaultBranch.SHA),
		},
	})

	if err != nil {
		return g.handleCreateRefError(err)
	}

	return nil
}

func (g GithubRepositoryService) handleCreateRefError(err error) errors.SlugError {
	var errResponse *github.ErrorResponse
	if errors.As(err, &errResponse) {
		switch errResponse.Response.StatusCode {
		case http.StatusNotFound:
			return githubDomain.RepositoryNotFound().WrapError(err)
		case http.StatusConflict:
			return githubDomain.RepositoryAccessDenied().WrapError(err)
		}
	}

	return githubDomain.UnknownError().WrapError(err)
}

func (g GithubRepositoryService) handleRepositoryError(err error) errors.SlugError {
	var errResponse *github.ErrorResponse
	if errors.As(err, &errResponse) {
		switch errResponse.Response.StatusCode {
		case http.StatusNotFound:
			return githubDomain.RepositoryNotFound().WrapError(err)
		case http.StatusForbidden:
			return githubDomain.RepositoryAccessDenied().WrapError(err)
		}
	}

	return githubDomain.UnknownError().WrapError(err)
}

func (g GithubRepositoryService) handleRefError(err error) errors.SlugError {
	var errResponse *github.ErrorResponse
	if errors.As(err, &errResponse) && errResponse.Response.StatusCode == http.StatusNotFound {
		return githubDomain.RefNotFound().WrapError(err)
	}

	return githubDomain.UnknownError().WrapError(err)
}
