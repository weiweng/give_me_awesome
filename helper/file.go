package helper

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func getExecRootPath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	exPath := filepath.Dir(ex)
	return exPath, nil
}

func GetBookList() (string, []string, error) {
	rootPath, err := getExecRootPath()
	if err != nil {
		return "", nil, err
	}
	bookPath := path.Join(rootPath, "./books")
	files, _ := ioutil.ReadDir(bookPath)

	ret := make([]string, 0, 4)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.HasSuffix(f.Name(), ".txt") {
			ret = append(ret, f.Name())
		}
	}
	return bookPath, ret, nil
}
