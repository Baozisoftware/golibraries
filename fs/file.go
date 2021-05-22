package fs

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/baozisoftware/golibraries/utils"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func CreateFile(fp string) (file *os.File, err error) {
	dir, _ := filepath.Split(fp)
	if dir != "" {
		err = CreateDir(dir)
		if err != nil {
			return
		}
	}
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

func ReadFileAllBytes(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err == nil {
		defer func() { _ = file.Close() }()
		return ioutil.ReadAll(file)
	}
	return []byte{}, err
}

func ReadFileAllString(filepath string) (string, error) {
	data, err := ReadFileAllBytes(filepath)
	if err == nil {
		return string(data), nil
	}
	return "", err
}

func ReadFileAllLines(filepath string) (lines []string, err error) {
	str, err := ReadFileAllString(filepath)
	if err == nil {
		s := "\n"
		if strings.Contains(str, "\r\n") {
			s = "\r\n"
		} else if strings.Contains(str, "\r") {
			s = "\r"
		}
		lines = strings.Split(str, s)
	}
	return
}

func WriteAllLinesToFile(fp string, lines []string) (err error) {
	s := utils.NewLine()
	data := strings.Join(lines, s)
	return WriteAllStringToFile(fp, data)
}

func WriteAllStringToFile(fp, str string) (err error) {
	return WriteAllBytesToFile(fp, []byte(str))
}

func WriteAllBytesToFile(fp string, data []byte) (err error) {
	dir, _ := filepath.Split(fp)
	if dir != "" {
		err = CreateDir(dir)
		if err != nil {
			return
		}
	}
	err = ioutil.WriteFile(fp, data, os.ModePerm)
	return
}
