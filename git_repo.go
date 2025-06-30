package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

type GitRepository struct {
	WorkTree string
	GitDir   string
	Config   *ini.File
}

func NewGitRepo(path string, force bool) (*GitRepository, error) {
	worktree := path
	gitdir := filepath.Join(path, ".git")

	if !force {
		info, err := os.Stat(gitdir)
		if err != nil || !info.IsDir() {
			return nil, fmt.Errorf("not a git repository: %s", path)
		}
	}
	repo := &GitRepository{
		WorkTree: worktree,
		GitDir:   gitdir,
	}
	configPath, err := repoFile(repo, false, "config")
	if err != nil {
		return nil, fmt.Errorf("")
	}
	if _, err := os.Stat(configPath); err == nil {
		config, err := ini.Load(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		repo.Config = config
	} else if !force {
		return nil, fmt.Errorf("configuration file missing")
	}
	if !force {
		vers, err := repo.Config.Section("core").Key("repositoryformatversion").Int()
		if err != nil {
			return nil, fmt.Errorf("missing or invalid repositoryformatversion")
		}
		if vers != 0 {
			return nil, fmt.Errorf("unsupported repositoryformatversion '%d': %w", vers, err)
		}
	}
	return repo, nil
}

func repoPath(repo *GitRepository, path ...string) string {
	return filepath.Join(append([]string{repo.GitDir}, path...)...)
}

// repoFile creates directories if it doesn't exist, excluding the last file.
//
// example: repoFile(repo, false, "refs", "remote", "origin", "HEAD") creates ".git/refs/remote/origin"
func repoFile(repo *GitRepository, mkdir bool, path ...string) (string, error) {
	if len(path) == 0 {
		return "", fmt.Errorf("no path given")
	}
	parentPath := path[:len(path)-1]

	_, err := repoDir(repo, mkdir, parentPath...)
	if err != nil {
		return "", err
	}
	return repoPath(repo, path...), nil
}

func repoDir(repo *GitRepository, mkdir bool, path ...string) (string, error) {
	dirPath := repoPath(repo, path...)

	if info, err := os.Stat(dirPath); err == nil {
		if info.IsDir() {
			return dirPath, nil
		}
		return "", fmt.Errorf("not a directory %s", dirPath)
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("error checking path %s: %w", dirPath, err)
	}
	if mkdir {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return "", fmt.Errorf("couldn't create dir %s: %w", dirPath, err)
		}
		return dirPath, nil
	}
	return "", nil
}

func repoCreate(path string) (*GitRepository, error) {
	repo, err := NewGitRepo(path, true)
	if err != nil {
		return nil, fmt.Errorf("error creating new repository: %w", err)
	}
	if info, err := os.Stat(repo.WorkTree); err == nil {
		if !info.IsDir() {
			return nil, fmt.Errorf("%s is not a directory", path)
		}
		entries, err := os.ReadDir(repo.GitDir)
		if err == nil && len(entries) > 0 {
			return nil, fmt.Errorf("%s is not empty", path)
		}
	} else if os.IsNotExist(err) {
		if err := os.MkdirAll(repo.WorkTree, 0755); err != nil {
			return nil, fmt.Errorf("failed to create worktree: %w", err)
		}
	} else {
		return nil, fmt.Errorf("error checking worktree: %w", err)
	}

	assertNoErr(repoDir(repo, true, "branches"))
	assertNoErr(repoDir(repo, true, "objects"))
	assertNoErr(repoDir(repo, true, "refs", "tags"))
	assertNoErr(repoDir(repo, true, "refs", "heads"))

	descPath, err := repoFile(repo, false, "description")
	if err != nil {
		return nil, err
	}
	content := []byte("Unnamed repository; edit this file 'description' to name the repository.\n")
	err = os.WriteFile(descPath, content, 0644)
	if err != nil {
		return nil, fmt.Errorf("error writing to file '%s': %w", descPath, err)
	}
	HEADPath, err := repoFile(repo, false, "HEAD")
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(HEADPath, []byte("ref: refs/head/master\n"), 0644)
	if err != nil {
		return nil, fmt.Errorf("error writing to file '%s': %w", descPath, err)
	}
	config, err := repoDefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("error initializing config: %w", err)
	}
	configPath, err := repoFile(repo, false, "config")
	if err != nil {
		return nil, err
	}
	err = config.SaveTo(configPath)
	if err != nil {
		return nil, fmt.Errorf("error writing to config file '%s': %w", configPath, err)
	}
	return repo, nil
}

func repoDefaultConfig() (*ini.File, error) {
	config := ini.Empty()
	core, err := config.NewSection("core")
	if err != nil {
		return nil, err
	}
	_, _ = core.NewKey("repositoryformatversion", "0")
	_, _ = core.NewKey("filemode", "false")
	_, _ = core.NewKey("bare", "false")
	return config, nil
}

func assertNoErr(_ any, err error) {
	if err != nil {
		panic(err)
	}
}
