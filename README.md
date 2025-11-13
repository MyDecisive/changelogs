# changelogs
Changelogs for MyDecisive components

### Table of Contents
- [Use Generate Changelog Reusable Workflow](#use-generate-changelog-reusable-workflow)
  - [Inputs](#inputs)
  - [Outputs](#outputs)
  - [Example Usage](#example-usage)

## Use Generate Changelog Reusable Workflow

This repository provides a reusable workflow that can be used to generate (using [git-cliff](https://git-cliff.org/)) and commit the changelog in the caller repository.

### Inputs
- `ref`: The branch, tag or SHA used to generate changelog (default: `main`)
- `config-url`: URL of the git cliff configuration file (default: `https://raw.githubusercontent.com/DecisiveAI/changelogs/refs/heads/main/cliff.toml`)
- `dry-run`: True to generate changelog without committing (default: `false`)

### Outputs
- `changelog`: Content of the generated changelog

### Example Usage

#### Basic

This example showcase how to use the workflow with default values:

```yaml
jobs:
  changelog:
    uses: DecisiveAI/changelogs/.github/workflows/reusable-changelog-gen.yaml@main
    secrets: inherit # pass all secrets
```

#### Dry Run

This example showcase how to run the workflow without making any commits:

```yaml
jobs:
  changelog:
    uses: DecisiveAI/changelogs/.github/workflows/reusable-changelog-gen.yaml@main
    with:
      dry-run: true
    secrets: inherit # pass all secrets
```