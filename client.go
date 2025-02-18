// This is the main package
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "app"}

	var initCommand = &cobra.Command{
		Use:   "init",
		Short: "Initialize a GoGit directory",
		Run: func(_ *cobra.Command, _ []string) {
			goInit()
		},
	}

	rootCmd.AddCommand(initCommand)
	rootCmd.Execute()
}

func goInit() {
	currDir, err := os.Getwd()
	if err != nil {
		panic(-1)
	}
	fmt.Printf("Initialized empty gogit repository in %s/%s", currDir, GitDir)
	Init()
}
