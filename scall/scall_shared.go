package scall

import (
	"github.com/Baozisoftware/golibraries/fs"
	"io"
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

func CreateProcess(stdin io.Reader, stdout, stderr io.Writer, prog string, args ...string) (p *os.Process, err error) {
	cmd := exec.Command(prog, args...)
	if stdin != nil {
		cmd.Stdin = stdin
	}
	if stdout != nil {
		cmd.Stdout = stdout
	}
	if stderr != nil {
		cmd.Stderr = stderr
	}
	err = cmd.Start()
	if err == nil {
		p = cmd.Process
	}
	return
}