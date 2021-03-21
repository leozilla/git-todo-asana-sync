package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os/exec"
)

func RunCommandWithOutput(cmd *exec.Cmd) ([]byte, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	errData, err := ioutil.ReadAll(stderr)
	if err != nil {
		return nil, err
	}

	err = cmd.Wait()
	if err != nil {
		return nil, err
	}

	outData, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}

	if len(errData) > 0 {
		return nil, errors.New(fmt.Sprintf("Executing command %s was not successful. Failed with %s", cmd.String(), string(errData)))
	}

	return outData, nil
}
