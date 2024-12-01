# github batch updater
This application is a rough tool to help with mundane tasks like updating same GH workflow in multiple repos.

It is able to fetch files from multiple repositories, replace some strings and create a PR with the changes.
You can also use `yq` expressions to modify the yaml files.

## Installation

```bash
go install github.com/kolah/github-batch-updater/cmd/gbu@latest
```
## Usage

In order to interact with the GitHub API, you need to set the `GITHUB_PERSONAL_ACCESS_TOKEN` env variable with a personal access token.
The token should have the `repo` scope. If you want to interact with GitHub workflows, you need to add the `workflow` scope.

```bash
gbu run script.yml
```

## Example script
```yaml
create_pr: true
target_branch: batch-test
repositories:
  - name: github-batch-updater
    owner: kolah
files:
  .github/workflows/build_image.yaml:
    - replace:
        replaces:
          - match: "uses: examplecom/pipelines/.github/workflows/someworkflow.yaml@v1.3"
            replace: "uses: examplecom/pipelines/.github/workflows/someworkflow.yaml@v1.6.7"
    - yq:
        expressions:
          - 'del(.jobs.success)'
          - 'del(.jobs.error)'
          - 'del(.jobs.* | select(.uses == "owner/repo/.github/workflows/some_workflow.yaml*"))'
          - ".run-name = \"${{ inputs.revision == 'HEAD' && format('{0}: revision {1}', inputs.environment_name, github.sha) || format('{0}: revision {1}', inputs.environment_name, inputs.revision) }}\""
pull_request:
  title: "Automatic PR batch update"
  body: "This PR updates the build image to the latest version."
  reviewers:
    - "kolah"
  team_reviewers:
    - "my-team"
```
