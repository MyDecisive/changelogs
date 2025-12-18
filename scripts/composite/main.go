package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

// Please take a look at README.md to learn more about using this tool.

const (
	latestReleaseURLFormat = "https://api.github.com/repos/%s/releases/latest"
	tagsURLFormat          = "https://api.github.com/repos/%s/tags"
	fileURLFormat          = "https://raw.githubusercontent.com/%s/refs/tags/%s/Chart.yaml"
	cloneURLFormat         = "https://github.com/%s/%s.git"
)

const (
	latestTagKey        = "tag_name"
	tagNameKey          = "name"
	prereleaseDelimiter = "-"
)

// chart used to parse Chart.yaml to get dependencies.
type chart struct {
	Dependencies []chartDependency `yaml:"dependencies"`
}

// chartDependency represents the fields needed from dependencies property.
type chartDependency struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

func main() {
	gitOwner := flag.String("owner", "DecisiveAI", "GitHub owner of the repositories")
	gitMainRepo := flag.String("repo", "mdai-hub", "GitHub Helm chart repository to gather the dependencies from")
	version := flag.String("version", "", "GitHub Helm chart repository's release version to generate changelog for")
	identifier := flag.String("id", "mdai", "Identifier used to find relevant dependencies")
	config := flag.String("config", ".https://raw.githubusercontent.com/DecisiveAI/changelogs/refs/heads/main/scripts/composite/cliff.toml", "url of the cliff.toml to use")
	path := flag.String("path", "./../../CHANGELOG.md", "absolute path to store the composite changelog")
	flag.Parse()

	var err error

	helmRepo := fmt.Sprintf("%s/%s", *gitOwner, *gitMainRepo)

	latestTag := *version
	if latestTag == "" {
		latestTag, err = getLatestTag(helmRepo)
		if err != nil {
			log.Fatalf("Unable to get latest tag: %v", err)
		}
	}

	latestDep, err := getDependencies(*identifier, helmRepo, latestTag)
	if err != nil {
		log.Fatalf("Unable to get latest dependencies: %v", err)
	}

	previousDep, err := getPreviousDependencies(*identifier, helmRepo, latestTag)
	if err != nil {
		log.Fatalf("Unable to get previous dependencies: %v", err)
	}

	composite, err := genComposite(context.Background(), *gitOwner, *config, latestTag, latestDep, previousDep)
	if err != nil {
		log.Fatalf("Unable to generate composite changelog: %v", err)
	}

	if err = writeToFile(*path, composite); err != nil {
		log.Fatalf("Unable write to file: %v", err)
	}
}

// writeToFile prepends the given composite changelog to the file at given path.
func writeToFile(path string, composite []byte) error {
	// Prepend is not a functionality natively supported by OS' file system.
	//
	// So in order to perform prepend:
	// We must first write the data we want to prepend to a temporary file
	temp, err := os.CreateTemp(filepath.Dir(path), "changelog")
	if err != nil {
		return fmt.Errorf("create temporary file:%w", err)
	}
	defer os.Remove(temp.Name())
	temp.Write(composite)

	// Then append the original content at given path to the temporary file
	file, err := os.Open(path)
	if err != nil {
		// File doesn't exist, just rename the temporary file to given path
		if os.IsNotExist(err) {
			err = os.Rename(temp.Name(), path)
			if err != nil {
				return fmt.Errorf("rename to %s:%w", path, err)
			}
			return nil
		}
		return fmt.Errorf("open %s:%w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Fprintln(temp, scanner.Text())
	}
	temp.Sync()

	// And lastly, rename the temporary file to given path
	err = os.Rename(temp.Name(), path)
	if err != nil {
		return fmt.Errorf("write to %s:%w", path, err)
	}

	return nil
}

// genComposite generate composite changelog.
func genComposite(ctx context.Context, gitOwner string, config string, latestTag string, latestDep map[string]string, previousDep map[string]string) ([]byte, error) {
	// Store all the cloned repo in temporary directory
	directory, err := os.MkdirTemp("", "deprepo")
	if err != nil {
		return nil, fmt.Errorf("create temporary directory:%w", err)
	}
	defer os.RemoveAll(directory)

	result := []byte{}
	result = append(result, fmt.Sprintf("## %s\n", latestTag)...)
	for repo, latest := range latestDep {
		path := filepath.Join(directory, repo)
		// clone without checkout to allow git cliff to see the tag&commit info
		if err = gitClone(ctx, gitOwner, repo, path); err != nil {
			return nil, fmt.Errorf("get tag info for %s:%w", repo, err)
		}

		previous := previousDep[repo]
		if previous != "" {
			previous = fmt.Sprintf("v%s", previous)
		}

		changelog, err := gitCliff(ctx, config, path, previous, fmt.Sprintf("v%s", latest))
		if err != nil {
			return nil, fmt.Errorf("get changelog for %s:%w", repo, err)
		}

		// If changelog length is 2 then it only contains `\n`
		if len(changelog) > 2 {
			result = append(result, fmt.Sprintf("### %s\n", repo)...)
			result = append(result, changelog...)
		}
	}

	return result, nil
}

// getLatestTag return latest released tag of the given repo.
func getLatestTag(repo string) (string, error) {
	data, err := getHTTP(fmt.Sprintf(latestReleaseURLFormat, repo))
	if err != nil {
		return "", err
	}
	result := make(map[string]any)
	if err = json.Unmarshal(data, &result); err != nil {
		return "", err
	}

	tagRaw, ok := result[latestTagKey]
	if !ok {
		return "", errors.New("find latest tag")
	}

	tag, ok := tagRaw.(string)
	if !ok {
		return "", errors.New("cast latest tag to string")
	}

	return tag, nil
}

// getPreviousDependencies return relevant dependencies from the none prerelease tag before the given latest tag.
func getPreviousDependencies(identifier string, repo string, latest string) (map[string]string, error) {
	tags, err := getAllTag(repo)
	if err != nil {
		return nil, fmt.Errorf("get all tags:%w", err)
	}

	// Find previous tag index
	previousTagIndex := len(tags) + 1
	for i, tag := range tags {
		if tag == latest {
			previousTagIndex = i + 1
			break
		}
	}

	// Only get previous dependencies if previous tag exists
	var result map[string]string
	if previousTagIndex < len(tags) {
		previous := tags[previousTagIndex]
		result, err = getDependencies(identifier, repo, previous)
		if err != nil {
			return nil, fmt.Errorf("get dependencies:%w", err)
		}
	}

	return result, nil
}

// getAllTag get all none prerelease tags, in descending order, for the given repo.
func getAllTag(repo string) ([]string, error) {
	data, err := getHTTP(fmt.Sprintf(tagsURLFormat, repo))
	if err != nil {
		return nil, err
	}
	tags := []map[string]any{}
	if err = json.Unmarshal(data, &tags); err != nil {
		return nil, err
	}

	// Parse and remove prerelease tags
	result := []string{}
	for _, tag := range tags {
		raw, ok := tag[tagNameKey]
		if !ok {
			return nil, errors.New("name not found")
		}

		tagStr, ok := raw.(string)
		if !ok {
			return nil, errors.New("cast tag to string")
		}
		if strings.Contains(tagStr, prereleaseDelimiter) {
			continue
		}
		result = append(result, tagStr)
	}

	// Sort in descending order
	slices.Sort(result)
	slices.Reverse(result)
	return result, nil
}

// getDependencies get relevant dependencies' name and version from Helm Chart.yaml file.
// Note: Only dependencies that contain `identifier` will be returned.
func getDependencies(identifier string, repo string, tag string) (map[string]string, error) {
	data, err := getHTTP(fmt.Sprintf(fileURLFormat, repo, tag))
	if err != nil {
		return nil, err
	}

	dependencies := chart{}
	yaml.Unmarshal(data, &dependencies)

	result := make(map[string]string)
	for _, dep := range dependencies.Dependencies {
		if strings.Contains(dep.Name, identifier) {
			result[dep.Name] = dep.Version
		}
	}

	return result, nil
}

// gitClone executes git clone no checkout for the given repo and store the repo at provided path.
func gitClone(ctx context.Context, owner string, repo string, path string) error {
	return exec.CommandContext(ctx, "git", "clone", "--no-checkout", fmt.Sprintf(cloneURLFormat, owner, repo), path).Run()
}

// gitCliff executes git cliff command using given parameters and return output.
func gitCliff(ctx context.Context, config string, path string, prev string, latest string) ([]byte, error) {
	return exec.CommandContext(ctx, "git-cliff", "--config-url", config, "--workdir", path, fmt.Sprintf("%s..%s", prev, latest)).Output()
}

// getHTTP get content from given url.
func getHTTP(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return io.ReadAll(res.Body)
}
