package main

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
