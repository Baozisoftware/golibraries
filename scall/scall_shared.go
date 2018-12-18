package scall

import (
	"github.com/Baozisoftware/golibraries/fs"
	"os"
	"os/exec"
)

func GetExecutable() (full, dir, name, ext, namewithoutext string) {
	p, err := os.Executable()
	if err == nil {
		full = p
		dir, name, ext, namewithoutext = fs.SplitFileName(p)
	}
	return
}

func CreateProcess(prog string, args ...string) (p *os.Process, err error) {
	cmd := exec.Command(prog, args...)
	err = cmd.Start()
	if err == nil {
		p = cmd.Process
	}
	return
}
