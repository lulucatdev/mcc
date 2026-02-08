//go:build !windows

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func launchClaude(profilePath string, extraEnv []string) error {
	// Find claude executable
	claudePath, err := exec.LookPath("claude")
	if err != nil {
		return fmt.Errorf("claude not found in PATH: %w", err)
	}

	// Set CLAUDE_CONFIG_DIR to the actual profile directory (not the symlink)
	// so that concurrent instances each use their own profile
	env := os.Environ()
	env = append(env, fmt.Sprintf("CLAUDE_CONFIG_DIR=%s", profilePath))
	env = append(env, extraEnv...)

	// Use syscall.Exec to replace current process with claude
	return syscall.Exec(claudePath, []string{"claude"}, env)
}
