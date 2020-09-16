package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func DockerHostConfigBinds(source, target string) string {
	return fmt.Sprintf("%s:%s", source, target)
}

func SubmissionBaseDir() (string, error) {
	return filepath.Abs("./files/submission/")
}

func SubmissionGenerateDirWithMkdir(sid uint32) (string, error) {
	t := time.Now()

	path, err := SubmissionBaseDir()
	if err != nil {
		return "", err
	}
	path, err = filepath.Abs(fmt.Sprintf("%s/%s/%d", path, t.Format("2006/01/02"), sid))
	if err != nil {
		return "", err
	}

	if !Exists(path) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return "", err
		}
	}

	return path, nil
}

func JudgeGenerateOutputDirWithMkdir(baseDir string) (string, error) {
	path := baseDir + "/output"

	if !Exists(path) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return "", err
		}
	}

	return path, nil
}

func JudgeCaseDir(tid uint32, version string) (string, error) {
	return filepath.Abs(fmt.Sprintf("./files/judge/%d/%s", tid, version))
}

func JudgeFilePath(tid uint32, version, caseId, caseType string) (string, error) {
	return filepath.Abs(fmt.Sprintf("./files/judge/%d/%s/%s.%s", tid, version, caseId, caseType))
}
