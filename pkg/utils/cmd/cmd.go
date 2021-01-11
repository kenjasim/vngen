package utils

import (
	"bytes"
	"os/exec"

	"github.com/pkg/errors"
)

// Run a command in the command line
func Run(name string, arg ...string) (out string, err error) {

	//Set the command and run it
	cmd := exec.Command(name, arg...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	stdout, err := cmd.Output()

	if err != nil {
		return string(stdout), errors.Wrap(err, "failed on Run:")
	}

	return string(stdout), nil
}

// Output - Run a command in the command line and get stdout, stderr and golang err
func Output(name string, arg ...string) (stdout, stderr string, err error) {

	//Set the command and run it
	cmd := exec.Command(name, arg...)
	var stderrBuff bytes.Buffer
	cmd.Stderr = &stderrBuff
	stdoutByte, err := cmd.Output()

	if err != nil {
		return string(stdoutByte), string(stderrBuff.Bytes()), errors.Wrap(err, "failed on Output:")
	}

	return string(stdoutByte), string(stderrBuff.Bytes()), nil
}
