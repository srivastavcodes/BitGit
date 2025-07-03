package main

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

type GitObject interface {
	Serialize(repo *GitRepository) ([]byte, error)
	Deserialize(data []byte) error
}

func NewGitObject(data []byte) (GitObject, error) {
	return nil, errors.New("unimplemented: concrete GitObject type needed")
}

func objectRead(repo *GitRepository, sha string) (GitObject, error) {
	objPath, err := repoFile(repo, false, "objects", sha[0:2], sha[2:])
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(objPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exists: %w", err)
	}

	file, err := os.Open(objPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader, err := zlib.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("zlib decompression failed: %w", err)
	}
	defer reader.Close()

	raw, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("couldn't read decompressed data: %w", err)
	}

	space := bytes.IndexByte(raw, ' ')
	if space == -1 {
		return nil, fmt.Errorf("malformed object %s: no space separator", sha)
	}
	oType := string(raw[0:space])

	x00 := bytes.IndexByte(raw[space:], '\x00')
	if x00 == -1 {
		return nil, fmt.Errorf("malformed object %s: no null separator", sha)
	}

	size, err := strconv.Atoi(string(raw[space+1 : x00]))
	// x00 is an index to subtracting len(raw) till index+1 gives content size
	if err != nil || size != len(raw)-x00-1 {
		return nil, fmt.Errorf("malformed object %s: invalid size", sha)
	}

	var object GitObject
	switch oType {
	case "blob":
		object = &GitBlob{}
	default:
		return nil, fmt.Errorf("unknown type %s for object %s", oType, sha)
	}

	if err = object.Deserialize(raw[x00+1:]); err != nil {
		return nil, fmt.Errorf("error deserializing: %w", err)
	}
	return object, nil
}
