package scall

import (
	"os/exec"
	"os"
	"path/filepath"
	"strings"
	"io"
	"errors"
	"fmt"
	"crypto/md5"
	"encoding/hex"
)

func CreateProcess(prog string, args ...string) (p *os.Process, err error) {
	cmd := exec.Command(prog, args...)
	err = cmd.Start()
	if err == nil {
		p = cmd.Process
	}
	return
}

func CreateDir(dirpath string) error {
	return os.MkdirAll(dirpath, os.ModePerm)
}

func CreateFile(fp string) (file *os.File, err error) {
	dir, _ := filepath.Split(fp)
	err = CreateDir(dir)
	if err == nil {
		file, err = os.Create(fp)
	}
	return
}

func OpenFile(filepath string) (file *os.File, err error) {
	file, err = os.Open(filepath)
	if err != nil {
		file, err = CreateFile(filepath)
	}
	return
}

func SplitFileName(p string) (dir, name, ext, namewithoutext string) {
	dir, name = filepath.Split(p)
	ext = filepath.Ext(name)
	n := strings.LastIndex(name, ".")
	if n > 0 {
		namewithoutext = name[:n]
	}
	return
}

func GetExecutable() (full, dir, name, ext, namewithoutext string) {
	p, err := os.Executable()
	if err == nil {
		full = p
		dir, name, ext, namewithoutext = SplitFileName(p)
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

func copyFile(src, dst string) error {
	e, f := FileOrFolderExists(src)
	if !e {
		return errors.New("src is not exists")
	}
	if f {
		return errors.New("src is not file")
	}
	if _, n, _, _ := SplitFileName(dst); n == "" {
		_, n, _, _ = SplitFileName(src)
		dst = fmt.Sprintf("%s/%s", dst, n)
	}
	sf, err := OpenFile(src)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := OpenFile(dst)
	if err != nil {
		return err
	}
	defer df.Close()

	_, err = io.Copy(df, sf)
	return err
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
	err := CopyFileOrFolder(src, dst)
	if err != nil {
		return err
	}
	return os.RemoveAll(src)
}

func GetFileMD5(path string) (string, error) {
	e, f := FileOrFolderExists(path)
	if !e {
		return "", errors.New("path is not exists")
	}
	if f {
		return "", errors.New("path is not file")
	}
	file, err := OpenFile(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := md5.New()
	buf := make([]byte, md5.BlockSize<<20)
	for {
		if n, err := file.Read(buf); err == nil {
			hash.Write(buf[:n])
		} else if n == 0 {
			break
		} else {
			return "", err
		}
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
