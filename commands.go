package main

import (
	"fmt"
	"os"
)

func CmdCatFile(objectName, objectType string) error {
	repo, err := repoFind(".", true)
	if err != nil {
		return err
	}
	return catFile(repo, objectName, objectType)
}

func CmdHashObjects(path, oType string, write bool) error {
	var repo *GitRepository
	var err error

	if write {
		repo, err = repoFind(".", true)
		if err != nil {
			return err
		}
	} else {
		repo = nil
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %s %w", path, err)
	}
	defer file.Close()

	sha, err := objectHash(file, oType, repo)
	if err != nil {
		return fmt.Errorf("failed to hash object: %w", err)
	}
	_, err = os.Stdout.WriteString(sha + "\n")
	return err
}
