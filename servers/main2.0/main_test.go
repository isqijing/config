package main

import (
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	count, err := CopyDir(".", "output")
	if err != nil {
		t.Log(err)
	}
	t.Log(count)
}

func CopyDir(src string, dst string) (count int, err error) {

	dirs, err := os.ReadDir(src)
	if err != nil {
		return 0, err
	}

	pathCursor := dst + "/" + src
	err = os.Mkdir(pathCursor, 0600)
	if err != nil {
		return 0, err
	}
	for _, dir := range dirs {
		if dir.IsDir() {
			countSub, err := CopyDir(src+"/"+dir.Name(), pathCursor)
			if err != nil {
				return 0, err
			}
			count += countSub
		}
		bytesFile, err := os.ReadFile(src + "/" + dir.Name())
		if err != nil {
			return 0, err
		}
		err = os.WriteFile(pathCursor+"/"+dir.Name(), bytesFile, 0600)
		if err != nil {
			return 0, err
		}
		count++
	}
	return count, nil
}
