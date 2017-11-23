package scall

import (
	"os/exec"
	"os"
	"path"
	"path/filepath"
	"strings"
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

func CreateFile(filepath string) (file *os.File, err error) {
	err = CreateDir(path.Dir(filepath))
	if err == nil {
		file, err = os.Create(filepath)
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
