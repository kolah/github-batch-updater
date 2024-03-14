package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v60/github"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"github.com/kolah/github-batch-updater/batch"
	"github.com/kolah/github-batch-updater/di"
	"github.com/spf13/cobra"
)

func runCommand(app *di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Run the application",
		Args:  cobra.MatchAll(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			fileName := args[0]
			c, err := app.LoadBatch(fileName)
			if err != nil {
				return err
			}

			for _, r := range c.Repositories {
				app.Logger().Info("Repository", "owner", r.Owner, "name", r.Name)

				repo, _, err := app.GithubRESTClient().Repositories.Get(cmd.Context(), r.Owner, r.Name)
				if err != nil {
					app.Logger().Error("Error while getting a repository", "error", err)
					continue
				}

				defaultBranch := repo.GetDefaultBranch()
				app.Logger().Info("Default branch", "name", defaultBranch)

				ref, _, err := app.GithubRESTClient().Git.GetRef(cmd.Context(), r.Owner, r.Name, "refs/heads/"+defaultBranch)
				if err != nil {
					app.Logger().Error("Error while getting a reference", "error", err)
					continue
				}

				mainSHA := *ref.Object.SHA

				type fileToUpdate struct {
					content string
					sha     string
				}
				modifiedFiles := make(map[string]fileToUpdate)
				for repoFileName, replaces := range c.Replaces {
					app.Logger().Info("File", "name", repoFileName)
					cnt, _, _, err := app.GithubRESTClient().Repositories.GetContents(cmd.Context(), r.Owner, r.Name, repoFileName, &github.RepositoryContentGetOptions{Ref: mainSHA})

					if err != nil {
						app.Logger().Error("Error while downloading a file", "error", err)
						continue
					}

					initialContent, err := cnt.GetContent()
					if err != nil {
						app.Logger().Error("Error while getting content", "error", err)
						continue
					}

					targetContent := initialContent

					for _, replace := range replaces {
						app.Logger().Info("Replace", "match", replace.Match, "replace", replace.Replace)
						targetContent = strings.ReplaceAll(targetContent, replace.Match, replace.Replace)
					}
					edits := myers.ComputeEdits(span.URIFromPath(repoFileName), initialContent, targetContent)

					if len(edits) == 0 {
						app.Logger().Info("No changes made in the file", "name", repoFileName)
						continue
					}

					diff := fmt.Sprint(gotextdiff.ToUnified(repoFileName, repoFileName, initialContent, edits))
					fmt.Println(diff)

					modifiedFiles[repoFileName] = fileToUpdate{
						content: targetContent,
						sha:     cnt.GetSHA(),
					}
				}

				if len(modifiedFiles) == 0 {
					app.Logger().Infof("No files were modified, skipping the repository %s/%s", r.Owner, r.Name)
					continue
				}

				if err = createBranch(cmd.Context(), c.BranchName, mainSHA, app, r); err != nil {
					app.Logger().Error("Error while creating a new branch", "error", err)

					continue
				}

				for modifiedFileName, content := range modifiedFiles {
					_, _, err := app.GithubRESTClient().Repositories.UpdateFile(cmd.Context(), r.Owner, r.Name, modifiedFileName, &github.RepositoryContentFileOptions{
						Message: github.String("chore: update " + modifiedFileName),
						Content: []byte(content.content),
						Branch:  &c.BranchName,
						SHA:     &content.sha,
					})

					if err != nil {
						app.Logger().Error("Error while updating a file", "error", err)
						continue
					}
				}

				if c.CreatePR {
					newPR := &github.NewPullRequest{
						Title:               github.String(c.PullRequest.Title),
						Head:                github.String(c.BranchName),
						Base:                github.String(defaultBranch),
						Body:                github.String(c.PullRequest.Body),
						MaintainerCanModify: github.Bool(true),
					}

					pr, _, err := app.GithubRESTClient().PullRequests.Create(cmd.Context(), r.Owner, r.Name, newPR)
					if err != nil {
						app.Logger().Error("Error while creating a pull request", "error", err)
						continue
					}

					app.Logger().Info("Pull request created", "url", pr.GetHTMLURL())

					if len(c.PullRequest.Reviewers) > 0 || len(c.PullRequest.TeamReviewers) > 0 {
						_, _, err = app.GithubRESTClient().PullRequests.RequestReviewers(cmd.Context(), r.Owner, r.Name, pr.GetNumber(), github.ReviewersRequest{
							Reviewers:     c.PullRequest.Reviewers,
							TeamReviewers: c.PullRequest.TeamReviewers,
						})

						if err != nil {
							app.Logger().Error("Error while requesting reviewers", "error", err)
						}
					}
				}
			}

			return nil
		},
	}

}

func createBranch(ctx context.Context, newBranchName, mainSHA string, app *di.Container, r batch.Repository) error {
	newRef := &github.Reference{Ref: github.String("refs/heads/" + newBranchName), Object: &github.GitObject{SHA: github.String(mainSHA)}}
	_, _, err := app.GithubRESTClient().Git.CreateRef(ctx, r.Owner, r.Name, newRef)

	return err
}
