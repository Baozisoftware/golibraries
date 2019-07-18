package fs

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func SplitFileName(p string) (dir, name, ext, namewithoutext string) {
	dir, name = filepath.Split(p)
	ext = filepath.Ext(name)
	n := strings.LastIndex(name, ".")
	if n > 0 {
		namewithoutext = name[:n]
	}
	return
}

func FileOrFolderExists(path string) (exists bool, isFolder bool) {
	f, err := os.Stat(path)
	exists = err == nil
	if exists {
		isFolder = f.IsDir()
	}
	return
}

func CopyFileOrFolder(src, dst string) error {
	e, f := FileOrFolderExists(src)
	if !e {
		return errors.New("src is not exists")
	}
	if !f {
		return copyFile(src, dst)
	}
	return copyFolder(src, dst)
}

func MoveFileOrFolder(src, dst string) error {
	return os.Rename(src, dst)
}
