package autodetect

import (
	"crypto/md5" // #nosec
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"fmt"
)

type buildToolInfo struct {
	globToDetect string
	tool         string
	injecter     Injecter
}

func DetectDirectoriesToCache(path string) error {
	var buildToolInfoMapping = []buildToolInfo{
		{
			globToDetect: "build.gradle",
			tool:         "gradle",
			injecter:     newGradleInjecter(),
		},
		{
			globToDetect: "WORKSPACE",
			tool:         "bazel",
			injecter:     newBazelInjecter(),
		},
	}

	for _, supportedTool := range buildToolInfoMapping {
		hash, dir, err := hashIfFileExist(path, supportedTool.globToDetect)
		if err != nil {
			return err
		}
		if hash == "" {
			hash, dir, err = hashIfFileExist(path, filepath.Join("**", supportedTool.globToDetect))
		}
		if err != nil {
			return err
		}
		if dir != "" && hash != "" {
			err = supportedTool.injecter.InjectTool()
			if err != nil {
				fmt.Printf("Error while auto-injecting for %s build tool: %s\n", supportedTool.tool, err.Error())
				continue
			}
		}
	}
	return nil
}

func hashIfFileExist(path, glob string) (string, string, error) {
	matches, _ := filepath.Glob(filepath.Join(path, glob))
	if len(matches) == 0 {
		return "", "", nil
	}

	return calculateMd5FromFiles(matches)
}

func calculateMd5FromFiles(fileList []string) (string, string, error) {
	rootMostFile := shortestPath(fileList)
	file, err := os.Open(rootMostFile)

	if err != nil {
		return "", "", err
	}

	dir, err := filepath.Abs(filepath.Dir(rootMostFile))

	if err != nil {
		return "", "", err
	}

	defer file.Close()

	if err != nil {
		return "", "", err
	}

	hash := md5.New() // #nosec
	_, err = io.Copy(hash, file)

	if err != nil {
		return "", "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), dir, nil
}

func shortestPath(input []string) (shortest string) {
	size := len(input[0])
	for _, v := range input {
		if len(v) <= size {
			shortest = v
			size = len(v)
		}
	}

	return
}
