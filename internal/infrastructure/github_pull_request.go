package infrastructure

import (
	"context"

	"github.com/google/go-github/v60/github"
	githubDomain "github.com/kolah/github-batch-updater/internal/domain/github"
)

type GithubPullRequestService struct {
	client *github.Client
}

func NewGithubPullRequestService(client *github.Client) *GithubPullRequestService {
	return &GithubPullRequestService{client: client}
}

func (g GithubPullRequestService) Create(ctx context.Context, repo githubDomain.Repository, targetBranch string, pullRequestTemplate githubDomain.PullRequestTemplate) (githubDomain.PullRequest, error) {
	newPR := &github.NewPullRequest{
		Title: github.String(pullRequestTemplate.Title),
		Body:  github.String(pullRequestTemplate.Body),
		Head:  github.String(targetBranch),
		Base:  github.String(repo.SourceBranch.Name),
	}
	pr, _, err := g.client.PullRequests.Create(ctx, repo.Owner, repo.Name, newPR)
	if err != nil {
		return githubDomain.PullRequest{}, err
	}

	return githubDomain.PullRequest{
		Number:  pr.GetNumber(),
		HTMLURL: pr.GetHTMLURL(),
	}, nil
}

func (g GithubPullRequestService) RequestReview(ctx context.Context, repo githubDomain.Repository, pr githubDomain.PullRequest, prTemplate githubDomain.PullRequestTemplate) error {
	_, _, err := g.client.PullRequests.RequestReviewers(ctx, repo.Owner, repo.Name, pr.Number, github.ReviewersRequest{
		Reviewers:     prTemplate.Reviewers,
		TeamReviewers: prTemplate.TeamReviewers,
	})

	return err
}
