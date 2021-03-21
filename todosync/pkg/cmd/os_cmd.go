package cmd

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"os/exec"
)

func RunCommand(cmd *exec.Cmd) ([]byte, error) {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	err = cmd.Wait()
	if err != nil {
		return nil, err
	}

	if len(stderr.Bytes()) > 0 {
		return nil, errors.New(fmt.Sprintf("Executing command %s was not successful. Failed with %s", cmd.String(), string(stderr.Bytes())))
	}

	return stdout.Bytes(), nil
}
