package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fernandoeho/smiddy/internal/fs"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove old sprint directories",
	Run: func(cmd *cobra.Command, args []string) {
		if !fs.FileExists(".smiddy/sprint") {
			fmt.Println("No sprint directory found. Nothing to clean.")
			return
		}

		n, err := fs.NextSprintNumber(".smiddy/sprint")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if n == 1 {
			fmt.Println("No sprints to clean.")
			return
		}

		current := n - 1
		fmt.Printf("This will delete all sprint directories except sprint %d.\n", current)
		fmt.Print("Are you sure? (yes/no) > ")

		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))

		if answer != "yes" {
			fmt.Println("Aborted.")
			return
		}

		for i := 1; i < current; i++ {
			path := fmt.Sprintf(".smiddy/sprint/%d", i)
			if fs.FileExists(path) {
				if err := os.RemoveAll(path); err != nil {
					fmt.Fprintf(os.Stderr, "could not remove %s: %v\n", path, err)
					continue
				}
				fmt.Printf("  removed  %s\n", path)
			}
		}

		fmt.Println("Done.")
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
