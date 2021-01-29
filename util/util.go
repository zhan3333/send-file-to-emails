package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func FindFileName1(name string, dir string) (string, error) {
	files, err := filepath.Glob(fmt.Sprintf("%s/*%s*", dir, name))
	if err != nil {
		return "", err
	}
	if len(files) != 1 {
		return "", fmt.Errorf("find %d file name like %s", len(files), name)
	}
	return files[0], nil
}

func FindFileName(name string, dir string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if strings.Contains(file.Name(), name) {
			return file.Name(), nil
		}
	}
	return "", fmt.Errorf("未找到文件")
}

func IsFileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
