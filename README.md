# changelogs
Changelogs for MyDecisive components

### Table of Contents
- [Use Semantic PR Title Reusable Workflow](#use-semantic-pr-title-reusable-workflow)
    - [Example Usage](#example-usage)

## Use Semantic PR Title Reusable Workflow

This repository provides a reusable workflow that can be used to lint (using [action-semantic-pull-request](https://github.com/amannn/action-semantic-pull-request)) the PR title according to [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/).

The workflow will comment on the PR if there are any errors and subsequently delete the comment if the check passes.

### Example Usage

#### Adding a new workflow specifically for the `pull_request` event

```yaml
name: PR Title Validate

on:
  pull_request:
    types:
      - opened
      - reopened
      - edited
      - synchronize # as a required action

permissions:
  pull-requests: write

jobs:
  reusable:
    uses: DecisiveAI/changelogs/.github/workflows/reusable-semantic-pr-title.yaml@main
    secrets: inherit # pass all secrets
```

#### Marking the job as required for the repo (one-time setup)
Add a branch ruleset targeting all branches, then check the `Require status checks to pass` checkbox, and select the `reusable/Validate PR title` check.
![img_1.png](docs/images/enabling_required_status_check.png)

Then, it will require a passing check before merging like the below screenshot.

![img_2.png](docs/images/required_check.png)

### For work-in-progress PRs
Refer to the [WIP](https://github.com/amannn/action-semantic-pull-request?tab=readme-ov-file#work-in-progress-pull-requests) documentation for the action if bypassing the lint is desired.
