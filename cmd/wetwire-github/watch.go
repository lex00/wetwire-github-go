package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch [path]",
	Short: "Watch Go source files and auto-rebuild workflows",
	Long: `Monitor Go source files for changes and automatically rebuild GitHub Actions
workflow YAML files when changes are detected.

The watch command uses file system notifications to detect changes to .go files
in the specified directory (and subdirectories). When changes are detected, it
automatically triggers a rebuild with debouncing to avoid excessive rebuilds.

Examples:
  wetwire-github watch ./workflows
  wetwire-github watch ./workflows --output .github/workflows
  wetwire-github watch ./workflows --debounce 500ms
  wetwire-github watch ./workflows --lint-only

The command will:
  - Monitor all .go files in the specified directory
  - Debounce rapid changes (default 300ms)
  - Show timestamps for each rebuild
  - Display clear success/failure messages
  - Continue watching until interrupted (Ctrl+C)`,
	Args: cobra.MaximumNArgs(1),
	RunE: runWatch,
}

func init() {
	watchCmd.Flags().StringP("output", "o", ".github/workflows", "Output directory for generated YAML")
	watchCmd.Flags().StringP("debounce", "d", "300ms", "Debounce duration for file changes")
	watchCmd.Flags().Bool("lint-only", false, "Only run lint on changes, don't rebuild")
}

func runWatch(cmd *cobra.Command, args []string) error {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	output, _ := cmd.Flags().GetString("output")
	debounceStr, _ := cmd.Flags().GetString("debounce")
	lintOnly, _ := cmd.Flags().GetBool("lint-only")

	// Parse debounce duration
	debounceDuration, err := time.ParseDuration(debounceStr)
	if err != nil {
		return fmt.Errorf("invalid debounce duration: %w", err)
	}

	// Convert path to absolute
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}

	// Verify path exists
	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("path not found: %s", path)
	}

	// Create file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("create watcher: %w", err)
	}
	defer watcher.Close()

	// Add path to watcher
	if err := addWatchPaths(watcher, absPath, info); err != nil {
		return fmt.Errorf("add watch paths: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "[%s] Watching %s for changes (debounce: %s)\n", formatTimestamp(time.Now()), absPath, debounceStr)
	fmt.Fprintln(cmd.OutOrStdout(), "Press Ctrl+C to stop watching")

	// Do initial build/lint
	if lintOnly {
		fmt.Fprintf(cmd.OutOrStdout(), "[%s] Initial lint...\n", formatTimestamp(time.Now()))
		if err := runWatchLint(absPath); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "[%s] Lint failed: %v\n", formatTimestamp(time.Now()), err)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "[%s] Lint passed\n", formatTimestamp(time.Now()))
		}
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "[%s] Initial build...\n", formatTimestamp(time.Now()))
		if err := runWatchBuild(absPath, output); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "[%s] Build failed: %v\n", formatTimestamp(time.Now()), err)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "[%s] Build successful: %s\n", formatTimestamp(time.Now()), output)
		}
	}

	// Start watching for changes
	debounceTimer := time.NewTimer(0)
	if !debounceTimer.Stop() {
		<-debounceTimer.C
	}

	pendingRebuild := false

	for {
		select {
		case <-cmd.Context().Done():
			return cmd.Context().Err()

		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			// Filter events
			if !shouldProcessEvent(event.Op.String(), event.Name) {
				continue
			}

			// Reset debounce timer
			pendingRebuild = true
			debounceTimer.Reset(debounceDuration)

		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "[%s] Watch error: %v\n", formatTimestamp(time.Now()), err)

		case <-debounceTimer.C:
			if pendingRebuild {
				pendingRebuild = false

				if lintOnly {
					fmt.Fprintf(cmd.OutOrStdout(), "[%s] Linting...\n", formatTimestamp(time.Now()))
					if err := runWatchLint(absPath); err != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "[%s] Lint failed: %v\n", formatTimestamp(time.Now()), err)
					} else {
						fmt.Fprintf(cmd.OutOrStdout(), "[%s] Lint passed\n", formatTimestamp(time.Now()))
					}
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "[%s] Rebuilding...\n", formatTimestamp(time.Now()))
					if err := runWatchBuild(absPath, output); err != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "[%s] Build failed: %v\n", formatTimestamp(time.Now()), err)
					} else {
						fmt.Fprintf(cmd.OutOrStdout(), "[%s] Build successful: %s\n", formatTimestamp(time.Now()), output)
					}
				}
			}
		}
	}
}

// addWatchPaths recursively adds directories to the watcher
func addWatchPaths(watcher *fsnotify.Watcher, path string, info os.FileInfo) error {
	if info.IsDir() {
		if err := watcher.Add(path); err != nil {
			return err
		}

		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				subPath := filepath.Join(path, entry.Name())
				// Skip hidden directories and vendor
				if strings.HasPrefix(entry.Name(), ".") || entry.Name() == "vendor" {
					continue
				}

				subInfo, err := entry.Info()
				if err != nil {
					continue
				}

				if err := addWatchPaths(watcher, subPath, subInfo); err != nil {
					return err
				}
			}
		}
	} else {
		dir := filepath.Dir(path)
		if err := watcher.Add(dir); err != nil {
			return err
		}
	}

	return nil
}

// isGoFile checks if a file has .go extension
func isGoFile(path string) bool {
	return strings.HasSuffix(path, ".go")
}

// shouldProcessEvent determines if an event should trigger a rebuild
func shouldProcessEvent(op, path string) bool {
	if !isGoFile(path) {
		return false
	}

	switch op {
	case "CREATE", "WRITE", "REMOVE", "RENAME":
		return true
	default:
		return false
	}
}

// formatTimestamp formats a time as HH:MM:SS
func formatTimestamp(t time.Time) string {
	return t.Format("15:04:05")
}

// runWatchBuild runs the build command for the given path
func runWatchBuild(sourcePath, outputDir string) error {
	result := runBuild(sourcePath, outputDir, false)
	if !result.Success {
		if len(result.Errors) > 0 {
			return fmt.Errorf("%s", result.Errors[0])
		}
		return fmt.Errorf("build failed")
	}
	return nil
}

// runWatchLint runs the lint command for the given path
func runWatchLint(sourcePath string) error {
	result := runLintPath(sourcePath)
	if !result.Success {
		if len(result.Errors) > 0 {
			return fmt.Errorf("%s", result.Errors[0])
		}
		return fmt.Errorf("lint failed")
	}
	return nil
}
