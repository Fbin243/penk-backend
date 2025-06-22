package gql

import (
	"fmt"
	"os"
	"os/exec"
)

func SetupBoilerplate(path string) error {
	// Get the current working directory (root)
	rootDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	// Create the directories if they don't exist
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Init module
	cmd := exec.Command("go", "mod", "init", "tenkhours/"+path)
	cmd.Dir = path
	_, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to init module: %v", err)
	}

	// Use the new module
	cmd = exec.Command("go", "work", "use", path)
	cmd.Dir = rootDir
	_, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to use the new module: %v", err)
	}

	// Create a `tools.go`
	toolsPath := path + "/tools.go"
	content := `//go:build tools
package tools

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/99designs/gqlgen/graphql/introspection"
)`

	err = os.WriteFile(toolsPath, []byte(content), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write tools.go: %v", err)
	}

	// Create `graph/`
	err = os.MkdirAll(path+"/graph", os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create graph/: %v", err)
	}

	// Install gqlgen v0.17.49
	cmd = exec.Command("go", "get", "-d", "github.com/99designs/gqlgen@v0.17.49")
	cmd.Dir = path
	_, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install gqlgen: %v", err)
	}

	// Build the server
	cmd = exec.Command("go", "run", "github.com/99designs/gqlgen", "init")
	cmd.Dir = path
	_, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to build the server: %v", err)
	}

	return nil
}
