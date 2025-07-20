package main

import "fmt"

type GitBlob struct {
	BlobData []byte
}

func (blob *GitBlob) Serialize(_ *GitRepository) ([]byte, error) {
	return blob.BlobData, nil
}

func (blob *GitBlob) Type() string { return "blob" }

func (blob *GitBlob) Deserialize(data []byte) error {
	blob.BlobData = data
	return nil
}

type GitTree struct{}

func (blob *GitTree) Serialize(_ *GitRepository) ([]byte, error) {
	return nil, fmt.Errorf("unimplemented object")
}

func (blob *GitTree) Type() string { return "tree" }

func (blob *GitTree) Deserialize(data []byte) error {
	return fmt.Errorf("unimplemented object")
}

type GitCommit struct{}

func (blob *GitCommit) Serialize(_ *GitRepository) ([]byte, error) {
	return nil, fmt.Errorf("unimplemented object")
}

func (blob *GitCommit) Type() string { return "commit" }

func (blob *GitCommit) Deserialize(data []byte) error {
	return fmt.Errorf("unimplemented object")
}

type GitTag struct{}

func (blob *GitTag) Serialize(_ *GitRepository) ([]byte, error) {
	return nil, fmt.Errorf("unimplemented object")
}

func (blob *GitTag) Type() string { return "tag" }

func (blob *GitTag) Deserialize(data []byte) error {
	return fmt.Errorf("unimplemented object")
}
