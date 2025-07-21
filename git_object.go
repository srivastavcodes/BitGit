package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"strconv"
)

type GitObject interface {
	Serialize(repo *GitRepository) ([]byte, error)
	Deserialize(data []byte) error
	Type() string
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

	x00 := bytes.IndexByte(raw, '\x00')
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
	case "tree":
		object = &GitTree{}
	case "tag":
		object = &GitTag{}
	case "commit":
		object = &GitCommit{}
	default:
		return nil, fmt.Errorf("unknown type %s for object %s", oType, sha)
	}

	if err = object.Deserialize(raw[x00+1:]); err != nil {
		return nil, fmt.Errorf("error deserializing: %w", err)
	}
	return object, nil
}

func objectWrite(object GitObject, repo *GitRepository) (string, error) {
	data, err := object.Serialize(repo)
	if err != nil {
		return "", fmt.Errorf("failed to serialize object: %w", err)
	}
	header := fmt.Sprintf("%s %d \x00", object.Type(), len(data))
	result := append([]byte(header), data...)

	// compute the hash
	sha := fmt.Sprintf("%x", sha1.Sum(result))

	if repo == nil {
		return sha, nil
	}
	// compute the path
	path, err := repoFile(repo, true, "objects", sha[0:2], sha[2:])
	if err != nil {
		return "", fmt.Errorf("failed to create object path: %w", err)
	}
	// check if the file already exists
	if _, err = os.Stat(path); os.IsNotExist(err) {
		var compressed bytes.Buffer
		writer := zlib.NewWriter(&compressed)

		if _, err := writer.Write(result); err != nil {
			return "", fmt.Errorf("failed to compressed data: %w", err)
		}
		if err = writer.Close(); err != nil {
			return "", fmt.Errorf("failed to close compressor: %w", err)
		}
		if err = os.WriteFile(path, compressed.Bytes(), 0644); err != nil {
			return "", fmt.Errorf("error writing objects to file: %w", err)
		}
	}
	return sha, nil
}

func catFile(repo *GitRepository, object, oType string) error {
	sha, err := objectFind(repo, object, oType, true)
	if err != nil {
		return err
	}
	objectData, err := objectRead(repo, sha)
	if err != nil {
		return err
	}
	data, err := objectData.Serialize(repo)
	if err != nil {
		return fmt.Errorf("failed to serialize repository: %w", err)
	}
	_, err = os.Stdout.Write(data)
	return err
}

func objectFind(_ *GitRepository, object string, _ string, _ bool) (string, error) {
	return object, nil
}

func objectHash(file io.Reader, oType string, repo *GitRepository) (string, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error while reading file: %w", err)
	}
	var object GitObject

	switch oType {
	case "blob":
		object = &GitBlob{}
	case "tree":
		object = &GitTree{}
	case "tag":
		object = &GitTag{}
	case "commit":
		object = &GitCommit{}
	default:
		return "", fmt.Errorf("unknown types %s", oType)
	}
	if err := object.Deserialize(data); err != nil {
		return "", fmt.Errorf("error while deserializing object: %w", err)
	}
	return objectWrite(object, repo)
}
