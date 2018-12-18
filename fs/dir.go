package fs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func CreateDir(dirpath string) error {
	return os.MkdirAll(dirpath, os.ModePerm)
}

func EachPath(path string, cb func(string) bool) error {
	return filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		if cb(path) {
			return nil
		}
		return errors.New("call cb return false")
	})
}

func FilterPath(path, expr string) (list []string, err error) {
	list = make([]string, 0)
	reg, e := regexp.Compile(expr)
	if e != nil {
		err = e
		return
	}
	err = EachPath(path, func(path string) bool {
		if expr == "" || reg.MatchString(path) {
			list = append(list, path)
		}
		return true
	})
	return
}

func copyFolder(src, dst string) error {
	e, f := FileOrFolderExists(src)
	if !e {
		return errors.New("src is not exists")
	}
	if !f {
		return errors.New("src is not folder")
	}
	if CreateDir(dst) != nil {
		return errors.New("faild to create dst folder")
	}
	s := len(src)
	if _, n, _, _ := SplitFileName(dst); n == "" {
		_, n, _, _ = SplitFileName(src)
		if n == "" {
			_, n, _, _ = SplitFileName(src[:len(src)-1])
		}
		dst = fmt.Sprintf("%s/%s", dst, n)
	}
	return filepath.Walk(src, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		return copyFile(path, dst+"/"+path[s:])
	})
}
