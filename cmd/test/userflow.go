package test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"tenkhours/cmd/auth"
)

func TestAPI(uid string, flows []string) error {
	// Get a new id token of the user
	token, err := auth.GetIdTokenByUID(uid)
	if err != nil {
		return fmt.Errorf("failed to get a new id token: %v", err)
	}

	rootDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	cmd := exec.Command("go", "test", "./test", "-v")
	cmd.Dir = rootDir
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("TOKEN=%s", token))
	cmd.Env = append(cmd.Env, fmt.Sprintf("FLOWS=%s", strings.Join(flows, ",")))

	// Redirect stdout and stderr to os.Stdout and os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %v", err)
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command finished with error: %v", err)
	}

	return nil
}
