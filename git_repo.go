package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type GitRepository struct {
	WorkTree string
	GitDir   string
	Config   string
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
	return repo, nil
}

func repoPath(repo *GitRepository, path ...string) string {
	return filepath.Join(append([]string{repo.GitDir}, path...)...)
}

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
