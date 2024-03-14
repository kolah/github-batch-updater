package batch

type Repository struct {
	Owner string
	Name  string
}

type FileReplace struct {
	Match   string
	Replace string
}

type PullRequest struct {
	Title         string
	Body          string
	Reviewers     []string
	TeamReviewers []string
}

type Config struct {
	Repositories []Repository
	Replaces     map[string][]FileReplace
	BranchName   string       `koanf:"target_branch_name"`
	CreatePR     bool         `koanf:"create_pr"`
	PullRequest  *PullRequest `koanf:"pull_request"`
}
