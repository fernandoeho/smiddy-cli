package cmd

import (
	"fmt"
	"os"

	"github.com/fernandoeho/smiddy/internal/fs"
	"github.com/fernandoeho/smiddy/internal/templates"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Smiddy project",
	Run: func(cmd *cobra.Command, args []string) {
		dirs := []string{
			".smiddy",
			".smiddy/sprints",
			".smiddy/agents",
		}

		for _, dir := range dirs {
			if err := fs.EnsureDir(dir); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		seeds := map[string]string{
			".smiddy/project-goals.md": templates.ProjectGoals,
			".smiddy/project-map.md":   templates.ProjectMap,
		}

		for path, content := range seeds {
			if fs.FileExists(path) {
				fmt.Printf("  skipped  %s (already exists)\n", path)
				continue
			}

			if err := fs.WriteFile(path, content); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Printf("  created  %s\n", path)
		}
		fmt.Println("\nProject initialized. Run `smiddy setup` to fill in your project goals.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
