package util

import (
	"fmt"
	"log"
	"os/exec"
)

func OsExec(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	stdoutStderr, err := cmd.CombinedOutput()
	fmt.Println(string(stdoutStderr))
	if err != nil {
		log.Println(err)
		return err
	}

	//fmt.Println(string(stdoutStderr))
	return nil
}
