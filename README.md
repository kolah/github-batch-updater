# github batch updater
This application is a rough tool to help with mundane tasks like updating same GH workflow in multiple repos.

It is able to fetch files from multiple repositories, replace some strings and create a PR with the changes.

## Installation

```bash
go install github.com/kolah/github-batch-updater/cmd/gbu@latest
```
## Usage

In order to interact with the GitHub API, you need to set the `GITHUB_PERSONAL_ACCESS_TOKEN` with a personal access token.
The token should have the `repo` scope. If you want to interact with github workflows, you need to add the `workflow` scope.

```bash
gbu run script.yml
```

## Example script
```yaml
batch:
  create_pr: true
  target_branch_name: batch-test
  repositories:
    - name: github-batch-updater
      owner: kolah
  replaces:
    .github/workflows/build_image.yaml:
      - match: "uses: examplecom/pipelines/.github/workflows/someworkflow.yaml@v1.3"
        replace: "uses: examplecom/pipelines/.github/workflows/someworkflow.yaml@v1.6.7"
      - match: "name: 'Weebhook for error'"
        replace: "name: 'Webhook for error'"
  pull_request:
    title: "Automatic PR batch update"
    body: "This PR updates the build image to the latest version."
    reviewers:
      - "kolah"
    team_reviewers:
      - "my-team"
```
