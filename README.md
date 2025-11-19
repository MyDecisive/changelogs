# changelogs
Changelogs for MyDecisive components

### Table of Contents
- [Use Semantic PR Title Reusable Workflow](#use-semantic-pr-title-reusable-workflow)
    - [Example Usage](#example-usage)

## Use Semantic PR Title Reusable Workflow

This repository provides a reusable workflow that can be used to lint (using [action-semantic-pull-request](https://github.com/amannn/action-semantic-pull-request)) the PR title according to [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/).

### Example Usage

#### Basic

```yaml
jobs:
  validate-pr-title:
    uses: DecisiveAI/changelogs/.github/workflows/reusable-semantic-pr-title.yaml@main
    secrets: inherit # pass all secrets
```

#### Requiring scope in the PR title

```yaml
jobs:
  validate-pr-title:
    uses: DecisiveAI/changelogs/.github/workflows/reusable-semantic-pr-title.yaml@main
    with:
      requireScope: true
    secrets: inherit # pass all secrets
```

### For work-in-progress PRs
Refer to the [WIP](https://github.com/amannn/action-semantic-pull-request?tab=readme-ov-file#work-in-progress-pull-requests) documentation for the action if bypassing the lint is desired.
