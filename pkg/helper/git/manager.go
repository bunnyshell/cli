package git

import (
	"net/url"
	"strconv"
	"strings"

	"bunnyshell.com/sdk"
)

type PrepareManager struct {
	Repositories []string

	repoDir map[string]string
	repoEnv map[string]string

	dirUsed map[string]bool
}

func NewPrepareManager() *PrepareManager {
	return &PrepareManager{
		Repositories: []string{},

		repoDir: map[string]string{},
		repoEnv: map[string]string{},

		dirUsed: map[string]bool{},
	}
}

func (pr *PrepareManager) AddComponents(components []sdk.ComponentGitCollection) error {
	for _, component := range components {
		if err := pr.addComponent(component); err != nil {
			return err
		}
	}

	return nil
}

func (pr *PrepareManager) GetDir(repository string) string {
	dir, exists := pr.repoDir[repository]
	if exists {
		return dir
	}

	panic("Missing call to either AddComponents")
}

func (pr *PrepareManager) GetEnvironment(repository string) string {
	environment, exists := pr.repoEnv[repository]
	if exists {
		return environment
	}

	panic("Missing call to either AddComponents or AddRepository")
}

func (pr *PrepareManager) hasDir(dir string) bool {
	return pr.dirUsed[dir]
}

func (pr *PrepareManager) hasRepo(repo string) bool {
	_, exists := pr.repoDir[repo]

	return exists
}

func (pr *PrepareManager) addComponent(component sdk.ComponentGitCollection) error {
	repository := component.GetRepository()
	if pr.hasRepo(repository) {
		return nil
	}

	dir, err := pr.makeDir(repository)
	if err != nil {
		return err
	}

	pr.dirUsed[dir] = true

	pr.Repositories = append(pr.Repositories, repository)
	pr.repoDir[repository] = dir
	pr.repoEnv[repository] = component.GetEnvironment()

	return nil
}

func (pr *PrepareManager) makeDir(repository string) (string, error) {
	parsed, err := url.Parse(repository)
	if err != nil {
		return "", err
	}

	parts := strings.Split(parsed.Path, "/")
	dirName := strings.TrimSuffix(parts[len(parts)-1], ".git")

	if !pr.hasDir(dirName) {
		return dirName, nil
	}

	suffix := 1

	for {
		ss := strconv.Itoa(suffix)
		if !pr.hasDir(dirName + ss) {
			return dirName + ss, nil
		}

		suffix++
	}
}
