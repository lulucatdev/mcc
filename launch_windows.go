//go:build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
)

func launchClaude(profilePath string, extraEnv []string) error {
	// Find claude executable
	claudePath, err := exec.LookPath("claude")
	if err != nil {
		return fmt.Errorf("claude not found in PATH: %w", err)
	}

	cmd := exec.Command(claudePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), fmt.Sprintf("CLAUDE_CONFIG_DIR=%s", profilePath))
	cmd.Env = append(cmd.Env, extraEnv...)

	return cmd.Run()
}
