package infrastructure

import (
	"context"
	"net/http"

	"github.com/google/go-github/v60/github"
	githubDomain "github.com/kolah/github-batch-updater/internal/domain/github"
	"github.com/kolah/github-batch-updater/internal/pkg/errors"
	"github.com/kolah/github-batch-updater/internal/pkg/ptr"
)

type GithubContentService struct {
	client *github.Client
}

func NewGithubContentService(client *github.Client) *GithubContentService {
	return &GithubContentService{client: client}
}

func (g GithubContentService) GetFile(ctx context.Context, repo githubDomain.Repository, path string) (*githubDomain.File, error) {
	contentResult, _, _, err := g.client.Repositories.GetContents(ctx, repo.Owner, repo.Name, path, &github.RepositoryContentGetOptions{
		Ref: repo.DefaultBranch.SHA,
	})
	if err != nil {
		return nil, g.handleGetFileError(err)
	}

	content, err := contentResult.GetContent()
	if err != nil {
		return nil, githubDomain.UnableToDecodeFileContent().WrapError(err)
	}

	return githubDomain.LoadGitHubFile(path, content, contentResult.GetSHA()), nil
}

func (g GithubContentService) UpdateContent(ctx context.Context, repo githubDomain.Repository, branch string, file *githubDomain.File) error {
	_, _, err := g.client.Repositories.UpdateFile(ctx, repo.Owner, repo.Name, file.Path(), &github.RepositoryContentFileOptions{
		Content: []byte(file.Content()),
		SHA:     ptr.To(file.SourceSHA()),
		Branch:  ptr.To(branch),
		Message: github.String("chore: " + file.Path()), // todo: commit message should be configurable
	})
	if err != nil {
		return githubDomain.FailedUpdatingFile().WrapError(err)
	}

	return nil
}

func (g GithubContentService) handleGetFileError(err error) errors.SlugError {
	var errResponse *github.ErrorResponse
	if errors.As(err, &errResponse) {
		switch errResponse.Response.StatusCode {
		case http.StatusNotFound:
			return githubDomain.FileNotFound().WrapError(err)
		case http.StatusForbidden:
			return githubDomain.FileAccessDenied().WrapError(err)
		}
	}

	return githubDomain.UnknownError().WrapError(err)
}
