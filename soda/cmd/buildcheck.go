package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
)

// checkBuildErrors runs a simple go build command to detect build errors in the project
// before attempting to run migrations. This prevents migration commands from hanging
// without feedback when there are build errors.
func checkBuildErrors() error {
	cmd := exec.Command("go", "build", "./...")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := stderr.String()
		if len(errMsg) > 0 {
			return fmt.Errorf("build errors detected before migration:\n%s", errMsg)
		}
		return fmt.Errorf("build command failed: %v", err)
	}

	return nil
}
