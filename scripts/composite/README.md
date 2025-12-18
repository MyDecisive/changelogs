# Generate Composite Changelog

Generate composite changelog base on relevant dependencies find in `Chart.yaml` of the given helm chart repository, and the generated composite changelog will be prepended to the given file path.

### Table of Contents
- [Installation](#installation)
- [Usage](#usage)
  - [Prerequisite](#prerequisite)
  - [Flags](#flags)
  - [Example(s)](#examples)
## Installation

```
go install github.com/DecisiveAI/changelogs/scripts/composite@latest
```

## Usage

### Prerequisite

Please make sure the following commands are installed and available in `$PATH` prior to running this tool:

- [git](https://git-scm.com/): Used to gather tag and commit information for the dependencies.
  - This tool only works with `v2.10.1`
- [git-cliff](https://git-cliff.org/docs/): Used to generate changelog for the dependencies.

### Flags

- `--owner`: GitHub owner of the repositories (default: `DecisiveAI`)
- `--repo`: GitHub Helm chart repository to gather dependencies from (default: `mdai-hub`)
- `--version`: GitHub Helm chart repository's release version to generate changelog for
  - If not provided, the tool will automatically generate changelog for the latest version
- `--id`: Identifier used to find relevant dependencies (default: `mdai`)
  - *Note: The dependencies find by the identifier must correspond to GitHub repositories with same name and be owned by the same GitHub owner as the one provided in `--owner`.*
    - *Version, without `v`, in the `Chart.yaml` must correspond to a tag, with `v`, in the corresponding GitHub repository.*
- `--config`: URL of the cliff.toml used to configure `git-cliff` (default: `./cliff.toml`)
- `--path`: Path to prepend the generated composite changelog to (default: `./../../CHANGELOG.md`)

### Example(s)

```
composite --config https://raw.githubusercontent.com/DecisiveAI/changelogs/refs/heads/main/scripts/composite/cliff.toml --path path/to/CHANGELOG.md
```
