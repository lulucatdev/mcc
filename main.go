package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	mccDirName       = ".mcc"
	profilesDirName  = "profiles"
	currentLinkName  = "current"
	configFileName   = "config.json"
	defaultProfile   = "default"
	profileMetaFile  = ".mcc-profile.json"
)

type Config struct {
	CurrentProfile string `json:"current_profile"`
}

type ProfileMeta struct {
	Provider string `json:"provider"`
	APIKey   string `json:"api_key"`
}

func loadProfileMeta(profilePath string) *ProfileMeta {
	data, err := os.ReadFile(filepath.Join(profilePath, profileMetaFile))
	if err != nil {
		return &ProfileMeta{Provider: "claude"}
	}
	var meta ProfileMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return &ProfileMeta{Provider: "claude"}
	}
	if meta.Provider == "" {
		meta.Provider = "claude"
	}
	return &meta
}

func saveProfileMeta(profilePath string, meta *ProfileMeta) error {
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(profilePath, profileMetaFile), data, 0644)
}

func getProviderEnv(meta *ProfileMeta) []string {
	switch meta.Provider {
	case "kimi":
		return []string{
			"ANTHROPIC_BASE_URL=https://api.kimi.com/coding/",
			"ANTHROPIC_API_KEY=" + meta.APIKey,
		}
	default: // "claude" or empty
		return nil
	}
}

func getMccDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}
	return filepath.Join(home, mccDirName)
}

func getProfilesDir() string {
	return filepath.Join(getMccDir(), profilesDirName)
}

func getCurrentLink() string {
	return filepath.Join(getMccDir(), currentLinkName)
}

func getConfigPath() string {
	return filepath.Join(getMccDir(), configFileName)
}

func getClaudeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}
	return filepath.Join(home, ".claude")
}

func loadConfig() (*Config, error) {
	configPath := getConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{CurrentProfile: defaultProfile}, nil
		}
		return nil, err
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func saveConfig(config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(getConfigPath(), data, 0644)
}

func profileExists(name string) bool {
	profilePath := filepath.Join(getProfilesDir(), name)
	info, err := os.Stat(profilePath)
	return err == nil && info.IsDir()
}

func listProfiles() ([]string, error) {
	profilesDir := getProfilesDir()
	entries, err := os.ReadDir(profilesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var profiles []string
	for _, entry := range entries {
		if entry.IsDir() {
			profiles = append(profiles, entry.Name())
		}
	}
	sort.Strings(profiles)
	return profiles, nil
}

func ensureMccStructure() error {
	profilesDir := getProfilesDir()

	// Create mcc directory
	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		return fmt.Errorf("failed to create profiles directory: %w", err)
	}

	// Check if default profile exists, if not, initialize from current .claude
	defaultProfileDir := filepath.Join(profilesDir, defaultProfile)
	if _, err := os.Stat(defaultProfileDir); os.IsNotExist(err) {
		claudeDir := getClaudeDir()
		if _, err := os.Stat(claudeDir); err == nil {
			// Copy current .claude to default profile
			if err := copyDir(claudeDir, defaultProfileDir); err != nil {
				return fmt.Errorf("failed to copy .claude to default profile: %w", err)
			}
			fmt.Println("Initialized default profile from existing ~/.claude")
		} else {
			// Create empty default profile
			if err := os.MkdirAll(defaultProfileDir, 0755); err != nil {
				return fmt.Errorf("failed to create default profile: %w", err)
			}
			fmt.Println("Created empty default profile")
		}
	}

	// Ensure current symlink exists
	currentLink := getCurrentLink()
	if _, err := os.Lstat(currentLink); os.IsNotExist(err) {
		if err := os.Symlink(defaultProfileDir, currentLink); err != nil {
			return fmt.Errorf("failed to create current symlink: %w", err)
		}
	}

	// Initialize config if not exists
	configPath := getConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := &Config{CurrentProfile: defaultProfile}
		if err := saveConfig(config); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
	}

	// Check and update shell config hint
	checkShellConfig()

	return nil
}

func checkShellConfig() {
	// Check if CLAUDE_CONFIG_DIR is already set correctly
	claudeConfigDir := os.Getenv("CLAUDE_CONFIG_DIR")
	expectedDir := getCurrentLink()

	if claudeConfigDir != expectedDir {
		fmt.Println("\n⚠️  To complete setup, add this to your ~/.zshrc or ~/.bashrc:")
		fmt.Printf("   export CLAUDE_CONFIG_DIR=\"%s\"\n", expectedDir)
		fmt.Println("   Then run: source ~/.zshrc (or restart your terminal)")
		fmt.Println()
	}
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy file
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dstPath, data, info.Mode())
	})
}

func copySettingsOnly(src, dst string) error {
	// List of settings files to copy (exclude credentials)
	settingsFiles := []string{
		"settings.json",
		"settings.local.json",
	}

	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	for _, filename := range settingsFiles {
		srcPath := filepath.Join(src, filename)
		if _, err := os.Stat(srcPath); err == nil {
			data, err := os.ReadFile(srcPath)
			if err != nil {
				continue
			}
			dstPath := filepath.Join(dst, filename)
			if err := os.WriteFile(dstPath, data, 0644); err != nil {
				return err
			}
		}
	}
	return nil
}

func switchProfile(name string, autoLaunch bool) error {
	if !profileExists(name) {
		return fmt.Errorf("profile '%s' does not exist. Use 'mcc new %s' to create it", name, name)
	}

	currentLink := getCurrentLink()
	profilePath := filepath.Join(getProfilesDir(), name)

	// Remove old symlink
	if err := os.Remove(currentLink); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove old symlink: %w", err)
	}

	// Create new symlink
	if err := os.Symlink(profilePath, currentLink); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	// Update config
	config := &Config{CurrentProfile: name}
	if err := saveConfig(config); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("✓ Switched to profile: %s\n", name)

	if autoLaunch {
		meta := loadProfileMeta(profilePath)
		extraEnv := getProviderEnv(meta)
		if meta.Provider != "claude" {
			fmt.Printf("  Launching claude (provider: %s)...\n", meta.Provider)
		} else {
			fmt.Println("  Launching claude...")
		}
		return launchClaude(profilePath, extraEnv)
	}
	return nil
}

func createProfile(name string, provider string, apiKey string) error {
	if profileExists(name) {
		return fmt.Errorf("profile '%s' already exists", name)
	}

	// Validate profile name
	if strings.ContainsAny(name, "/\\:*?\"<>|") {
		return fmt.Errorf("invalid profile name: contains forbidden characters")
	}

	profilePath := filepath.Join(getProfilesDir(), name)

	// Copy settings from default profile (without credentials)
	defaultProfileDir := filepath.Join(getProfilesDir(), defaultProfile)
	if err := copySettingsOnly(defaultProfileDir, profilePath); err != nil {
		// If copy fails, just create empty directory
		if err := os.MkdirAll(profilePath, 0755); err != nil {
			return fmt.Errorf("failed to create profile directory: %w", err)
		}
	}

	// Save profile metadata for non-claude providers
	if provider != "" && provider != "claude" {
		meta := &ProfileMeta{Provider: provider, APIKey: apiKey}
		if err := saveProfileMeta(profilePath, meta); err != nil {
			return fmt.Errorf("failed to save profile metadata: %w", err)
		}
	}

	fmt.Printf("✓ Created profile: %s\n", name)
	if provider != "" && provider != "claude" {
		fmt.Printf("  Provider: %s\n", provider)
	}
	fmt.Println()
	fmt.Println("To use this profile:")
	fmt.Printf("  mcc run %s\n", name)
	return nil
}

func deleteProfile(name string) error {
	if name == defaultProfile {
		return fmt.Errorf("cannot delete the default profile")
	}

	if !profileExists(name) {
		return fmt.Errorf("profile '%s' does not exist", name)
	}

	// Check if it's current profile
	config, err := loadConfig()
	if err != nil {
		return err
	}

	if config.CurrentProfile == name {
		return fmt.Errorf("cannot delete the currently active profile. Switch to another profile first")
	}

	profilePath := filepath.Join(getProfilesDir(), name)
	if err := os.RemoveAll(profilePath); err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}

	fmt.Printf("✓ Deleted profile: %s\n", name)
	return nil
}

func setAPIKey(name string, apiKey string) error {
	if !profileExists(name) {
		return fmt.Errorf("profile '%s' does not exist", name)
	}

	profilePath := filepath.Join(getProfilesDir(), name)
	meta := loadProfileMeta(profilePath)

	if meta.Provider == "claude" {
		return fmt.Errorf("profile '%s' uses the claude provider and does not need an API key", name)
	}

	meta.APIKey = apiKey
	if err := saveProfileMeta(profilePath, meta); err != nil {
		return fmt.Errorf("failed to save profile metadata: %w", err)
	}

	fmt.Printf("✓ Updated API key for profile: %s\n", name)
	return nil
}

func syncProfile(name string) error {
	if !profileExists(name) {
		return fmt.Errorf("profile '%s' does not exist. Use 'mcc new %s' to create it first", name, name)
	}

	claudeDir := getClaudeDir()
	info, err := os.Stat(claudeDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("~/.claude does not exist. Nothing to sync")
	}
	if err != nil {
		return fmt.Errorf("failed to access ~/.claude: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("~/.claude is not a directory")
	}

	// Check if ~/.claude is empty
	entries, err := os.ReadDir(claudeDir)
	if err != nil {
		return fmt.Errorf("failed to read ~/.claude: %w", err)
	}
	if len(entries) == 0 {
		return fmt.Errorf("~/.claude is empty. Nothing to sync")
	}

	profilePath := filepath.Join(getProfilesDir(), name)

	// Copy settings from ~/.claude (excluding credentials)
	count, skipped, err := syncSettings(claudeDir, profilePath)
	if err != nil {
		return fmt.Errorf("failed to sync settings: %w", err)
	}

	if count == 0 {
		fmt.Printf("⚠️  No settings files found in ~/.claude to sync\n")
		if skipped > 0 {
			fmt.Printf("   (%d credential file(s) were skipped)\n", skipped)
		}
		return nil
	}

	fmt.Printf("✓ Synced %d file(s) from ~/.claude to profile: %s\n", count, name)
	if skipped > 0 {
		fmt.Printf("  (%d credential file(s) were skipped for security)\n", skipped)
	}
	return nil
}

func syncSettings(src, dst string) (copied int, skipped int, err error) {
	// Files/patterns to exclude (credentials and auth-related)
	excludePatterns := []string{
		".credentials.json",
		"credentials.json",
		"auth.json",
		".auth",
	}

	// Directories to skip entirely
	skipDirs := []string{
		".git",
	}

	isExcluded := func(name string) bool {
		for _, pattern := range excludePatterns {
			if strings.Contains(strings.ToLower(name), strings.ToLower(pattern)) {
				return true
			}
		}
		return false
	}

	isSkipDir := func(name string) bool {
		for _, dir := range skipDirs {
			if name == dir {
				return true
			}
		}
		return false
	}

	err = filepath.Walk(src, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		// Skip certain directories entirely
		if info.IsDir() && isSkipDir(info.Name()) {
			return filepath.SkipDir
		}

		relPath, relErr := filepath.Rel(src, path)
		if relErr != nil {
			return relErr
		}

		// Skip excluded files (credentials)
		if isExcluded(info.Name()) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			skipped++
			return nil
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy file
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		if writeErr := os.WriteFile(dstPath, data, info.Mode()); writeErr != nil {
			return writeErr
		}
		copied++
		return nil
	})

	return copied, skipped, err
}

func showStatus() error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	profiles, err := listProfiles()
	if err != nil {
		return err
	}

	fmt.Println("Claude Code Account Manager (mcc)")
	fmt.Println()
	fmt.Printf("Current profile: %s\n", config.CurrentProfile)
	fmt.Println()
	fmt.Println("Available profiles:")

	for _, profile := range profiles {
		profilePath := filepath.Join(getProfilesDir(), profile)
		meta := loadProfileMeta(profilePath)
		providerTag := ""
		if meta.Provider != "claude" {
			providerTag = fmt.Sprintf(" [%s]", meta.Provider)
		}
		if profile == config.CurrentProfile {
			fmt.Printf("  * %s (active)%s\n", profile, providerTag)
		} else {
			fmt.Printf("    %s%s\n", profile, providerTag)
		}
	}

	// Check CLAUDE_CONFIG_DIR
	claudeConfigDir := os.Getenv("CLAUDE_CONFIG_DIR")
	expectedDir := getCurrentLink()

	fmt.Println()
	if claudeConfigDir == expectedDir {
		fmt.Println("✓ CLAUDE_CONFIG_DIR is correctly configured")
	} else if claudeConfigDir == "" {
		fmt.Println("⚠️  CLAUDE_CONFIG_DIR is not set")
		fmt.Printf("   Add to your shell config: export CLAUDE_CONFIG_DIR=\"%s\"\n", expectedDir)
	} else {
		fmt.Println("⚠️  CLAUDE_CONFIG_DIR points to a different location")
		fmt.Printf("   Current: %s\n", claudeConfigDir)
		fmt.Printf("   Expected: %s\n", expectedDir)
	}

	return nil
}

func showHelp() {
	fmt.Println("Claude Code Account Manager (mcc)")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  mcc                              Switch to default and launch claude")
	fmt.Println("  mcc run <name>                   Switch to profile and launch claude")
	fmt.Println("  mcc new <name>                   Create a new claude profile")
	fmt.Println("  mcc new <name> <provider> <key>  Create a profile with a provider")
	fmt.Println("  mcc set-key <name> <api-key>     Update API key for a profile")
	fmt.Println("  mcc sync [name]                  Sync ~/.claude to profile (default: current)")
	fmt.Println("  mcc status                       Show current status and profiles")
	fmt.Println("  mcc list                         List all profiles")
	fmt.Println("  mcc delete <name>                Delete a profile")
	fmt.Println("  mcc help                         Show this help message")
	fmt.Println()
	fmt.Println("Providers:")
	fmt.Println("  claude (default)  Standard Claude Code with Anthropic account")
	fmt.Println("  kimi              Kimi Coding (uses claude CLI with Kimi API)")
	fmt.Println()
	fmt.Println("Setup:")
	fmt.Println("  Add this to your ~/.zshrc or ~/.bashrc:")
	fmt.Printf("  export CLAUDE_CONFIG_DIR=\"%s\"\n", getCurrentLink())
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  mcc                              # Launch with default profile")
	fmt.Println("  mcc run work                     # Launch with 'work' profile")
	fmt.Println("  mcc new work                     # Create a claude profile")
	fmt.Println("  mcc new kimi-work kimi sk-xxx    # Create a Kimi profile")
	fmt.Println("  mcc set-key kimi-work sk-new     # Update API key")
	fmt.Println("  mcc status                       # Show all profiles")
}

func main() {
	// Ensure mcc structure exists
	if err := ensureMccStructure(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing mcc: %v\n", err)
		os.Exit(1)
	}

	args := os.Args[1:]

	// No args: switch to default and launch claude
	if len(args) == 0 {
		if err := switchProfile(defaultProfile, true); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	command := args[0]

	switch command {
	case "help", "-h", "--help":
		showHelp()

	case "status", "st":
		if err := showStatus(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "list", "ls":
		profiles, err := listProfiles()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing profiles: %v\n", err)
			os.Exit(1)
		}
		config, _ := loadConfig()
		for _, profile := range profiles {
			profilePath := filepath.Join(getProfilesDir(), profile)
			meta := loadProfileMeta(profilePath)
			providerTag := ""
			if meta.Provider != "claude" {
				providerTag = fmt.Sprintf(" [%s]", meta.Provider)
			}
			if profile == config.CurrentProfile {
				fmt.Printf("* %s%s\n", profile, providerTag)
			} else {
				fmt.Printf("  %s%s\n", profile, providerTag)
			}
		}

	case "new", "create", "add":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: profile name required")
			fmt.Fprintln(os.Stderr, "Usage: mcc new <name> [provider] [api-key]")
			os.Exit(1)
		}
		name := args[1]
		provider := ""
		apiKey := ""
		if len(args) >= 3 {
			provider = args[2]
		}
		if len(args) >= 4 {
			apiKey = args[3]
		}
		if provider != "" && provider != "claude" && apiKey == "" {
			fmt.Fprintf(os.Stderr, "Error: API key required for provider '%s'\n", provider)
			fmt.Fprintf(os.Stderr, "Usage: mcc new <name> %s <api-key>\n", provider)
			os.Exit(1)
		}
		if err := createProfile(name, provider, apiKey); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "delete", "rm", "remove":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: profile name required")
			fmt.Fprintln(os.Stderr, "Usage: mcc delete <name>")
			os.Exit(1)
		}
		name := args[1]
		if err := deleteProfile(name); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "sync":
		var name string
		if len(args) >= 2 {
			name = args[1]
		} else {
			// Use current profile
			config, err := loadConfig()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			name = config.CurrentProfile
		}
		if err := syncProfile(name); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "run":
		var name string
		if len(args) >= 2 {
			name = args[1]
		} else {
			name = defaultProfile
		}
		if err := switchProfile(name, true); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "set-key":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: profile name and API key required")
			fmt.Fprintln(os.Stderr, "Usage: mcc set-key <name> <api-key>")
			os.Exit(1)
		}
		name := args[1]
		apiKey := args[2]
		if err := setAPIKey(name, apiKey); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		fmt.Fprintln(os.Stderr, "Run 'mcc help' for usage")
		os.Exit(1)
	}
}
