package e2e

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"tenkhours/cmd/auth"
	"tenkhours/pkg/db"
	"tenkhours/services/core/repo"
)

func TestUserFlow(uid string) error {
	// Get a new id token of the user
	token, err := auth.GetIdTokenByUID(uid)
	if err != nil {
		return fmt.Errorf("failed to get a new id token: %v", err)
	}

	// Remove the current profile in `dev` database
	profileRepo := repo.NewProfilesRepo(db.GetDBManager().DB, nil)
	err = profileRepo.DeleteProfileByFirebaseUID(uid)
	if err != nil {
		return fmt.Errorf("failed to delete profile: %v", err)
	}

	// Remove the active session in redis
	redisClient := db.GetRedisClient()
	err = redisClient.Del(context.Background(), uid).Err()
	if err != nil {
		return fmt.Errorf("faild to delete the active session in redis")
	}

	rootDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	cmd := exec.Command("go", "test", "./test", "-v")
	cmd.Dir = rootDir
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("TOKEN=%s", token))

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
