package utils

import (
	"archive/tar"
	"bytes"
	"fmt"
	JudgerModels "github.com/Rabbit-OJ/Rabbit-OJ-Judger/models"
	"io"
	"io/ioutil"
	"path/filepath"
)

// this function is mainly modified from https://golang.org/pkg/archive/tar/

type TarFileBasicInfo struct {
	Name string
	Body []byte
	Mode int64
}

func ConvertToTar(files []TarFileBasicInfo) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	writer := tar.NewWriter(&buf)

	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Mode: file.Mode,
			Size: int64(len(file.Body)),
		}

		if err := writer.WriteHeader(hdr); err != nil {
			fmt.Println(err)
		}
		if _, err := writer.Write(file.Body); err != nil {
			fmt.Println(err)
		}
	}

	if err := writer.Close(); err != nil {
		fmt.Println(err)
	}
	return &buf, nil
}

func TarToFile(reader io.ReadCloser) ([]TarFileBasicInfo, error) {
	var files []TarFileBasicInfo

	tarReader := tar.NewReader(reader)
	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		fileBytes, err := ioutil.ReadAll(tarReader)
		if err != nil {
			return nil, err
		}

		files = append(files, TarFileBasicInfo{
			Name: hdr.FileInfo().Name(),
			Body: fileBytes,
			Mode: hdr.Mode,
		})
	}

	return files, nil
}

func AllFilesInDirToTarArchiveInfo(filePath, absPath string) ([]TarFileBasicInfo, error) {
	var basicInfo []TarFileBasicInfo

	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		name := file.Name()
		currentCasePath := filepath.Join(filePath, name)
		containerCasePath := filepath.Join(absPath, name)

		fileBytes, err := ioutil.ReadFile(currentCasePath)
		if err != nil {
			return nil, err
		}

		basicInfo = append(basicInfo, TarFileBasicInfo{
			Name: containerCasePath,
			Body: fileBytes,
			Mode: int64(file.Mode()),
		})
	}

	return basicInfo, nil
}

func TestCasesToTarArchiveInfo(testCases []*JudgerModels.TestCaseType, absPath string) ([]TarFileBasicInfo, error) {
	var basicInfo []TarFileBasicInfo

	for _, testCase := range testCases {
		id := testCase.Id
		fileName := fmt.Sprintf("%d.in", id)
		containerCasePath := filepath.Join(absPath, fileName)

		basicInfo = append(basicInfo, TarFileBasicInfo{
			Name: containerCasePath,
			Body: testCase.Stdin,
			Mode: 0644,
		})
	}

	return basicInfo, nil
}