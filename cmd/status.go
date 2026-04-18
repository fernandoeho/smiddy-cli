package cmd

import (
	"fmt"
	"os"

	"github.com/fernandoeho/smiddy/internal/fs"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current sprint status",
	Run: func(cmd *cobra.Command, args []string) {
		if !fs.FileExists(".smiddy") {
			fmt.Fprintln(os.Stderr, "No .smiddy directory found. Run `smiddy init` first.")
			os.Exit(1)
		}

		n, err := fs.NextSprintNumber(".smiddy/sprints")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if n == 1 {
			fmt.Println("No sprints yet. Run `smiddy new` to create one.")
			return
		}

		current := n - 1
		statusPath := fmt.Sprintf(".smiddy/sprint/%d/status.md", current)

		if !fs.FileExists(statusPath) {
			fmt.Printf("Sprint %d — no status file yet (not run)\n", current)
			return
		}

		content, err := fs.ReadFile(statusPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Printf("Sprint %d status:\n\n%s\n", current, content)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
