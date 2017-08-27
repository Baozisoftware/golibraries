package scall

import (
	"os/exec"
	"os"
	"path"
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
