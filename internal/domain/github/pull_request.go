package github

import (
	"context"
)

type PullRequestTemplate struct {
	Title         string
	Body          string
	Reviewers     []string
	TeamReviewers []string
}

type PullRequest struct {
	Number  int
	HTMLURL string
}

type PullRequestService interface {
	Create(ctx context.Context, repo Repository, targetBranch string, prTemplate PullRequestTemplate) (PullRequest, error)
	RequestReview(ctx context.Context, repo Repository, pr PullRequest, prTemplate PullRequestTemplate) error
}
