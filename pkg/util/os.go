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

func MergeMap(global, override map[string]interface{}) map[string]interface{} {
	if len(global) == 0 {
		return override
	}

	if len(override) == 0 {
		return global
	}

	for key, val := range override {
		global[key] = val
	}

	return global
}
