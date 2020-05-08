package loader

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Name      string
	Content   string
	Extension string
}

func getFile(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer closeFile(f)

	byteValue, err := ioutil.ReadAll(f)

	if err != nil {
		return nil, err
	}

	return &File{
		Name:      strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
		Content:   string(byteValue),
		Extension: filepath.Ext(path),
	}, nil
}

func listFilesByExtensions(basePath string, extensions ...string) ([]string, error) {
	files := []string{}

	if err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if isSupported := func(currentFileExtension string, supportedExtensions []string) bool {
			for _, e := range supportedExtensions {
				if strings.Trim(currentFileExtension, ".") == strings.Trim(e, ".") {
					return true
				}
			}

			return false
		}(filepath.Ext(path), extensions); isSupported == false {
			return nil
		}

		files = append(files, path)

		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}

func closeFile(file *os.File) {
	file.Close()
}
