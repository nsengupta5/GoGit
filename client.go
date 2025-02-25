// This is the main package
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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

	var hashObjectCommand = &cobra.Command{
		Use:   "hash-object",
		Short: "Hash the file contents of the given file",
		Run: func(_ *cobra.Command, args []string) {
			if len(args) > 1 {
				log.Fatal("Invalid number of arguments")
			}
			hashObject(args[0])
		},
	}

	var catFileCommand = &cobra.Command{
		Use:   "cat-file",
		Short: "Prints the content of an object",
		Run: func(_ *cobra.Command, args []string) {
			if len(args) > 1 {
				log.Fatal("Invalid number of arguments")
			}
			catFile(args[0])
		},
	}

	var writeTreeCommand = &cobra.Command{
		Use:   "write-tree",
		Short: "Writes the current working directory into object database",
		Run: func(_ *cobra.Command, args []string) {
			if len(args) > 1 {
				log.Fatal("Invalid number of arguments")
			}
			writeTree()
		},
	}

	var readTreeCommand = &cobra.Command{
		Use:   "read-tree",
		Short: "Reads the OID of a tree and writes the state to the current working directory",
		Run: func(_ *cobra.Command, args []string) {
			if len(args) > 1 {
				log.Fatal("Invalid number of arguments")
			}
			readTree(args[0])
		},
	}

	var commitMessage string
	var commitCommand = &cobra.Command{
		Use:   "commit",
		Short: "Creates a new commit",
		Run: func(_ *cobra.Command, args []string) {
			if commitMessage == "" {
				log.Fatal("Commit message is required. Use the -m flag to specify a commit message")
			}

			if len(args) > 1 {
				log.Fatal("Invalid number of arguments")
			}
			commit(commitMessage)
		},
	}
	commitCommand.Flags().StringVarP(&commitMessage, "message", "m", "", "commit message")

	var logCommand = &cobra.Command{
		Use:   "log",
		Short: "Display the logs of the commits",
		Run: func(_ *cobra.Command, _ []string) {
			goLog()
		},
	}

	var checkoutCommand = &cobra.Command{
		Use:   "checkout",
		Short: "Checkout to a specified commit",
		Run: func(_ *cobra.Command, args []string) {
			if len(args) > 1 {
				log.Fatal("Invalid number of arguments")
			}
			checkout(args[0])
		},
	}

	var tagCommand = &cobra.Command{
		Use:   "tag",
		Short: "Tag a specified commit with a specified name",
		Run: func(_ *cobra.Command, args []string) {
			if len(args) > 2 {
				log.Fatal("Invalid number of arguments")
			}
			tag(args[0], args[1])
		},
	}

	rootCmd.AddCommand(initCommand)
	rootCmd.AddCommand(hashObjectCommand)
	rootCmd.AddCommand(catFileCommand)
	rootCmd.AddCommand(writeTreeCommand)
	rootCmd.AddCommand(readTreeCommand)
	rootCmd.AddCommand(commitCommand)
	rootCmd.AddCommand(logCommand)
	rootCmd.AddCommand(checkoutCommand)
	rootCmd.AddCommand(tagCommand)
	rootCmd.Execute()
}

func goInit() {
	currDir, err := os.Getwd()
	if err != nil {
		panic(-1)
	}
	fmt.Printf("Initialized empty gogit repository in %s", filepath.Join(currDir, GitDir))
	Init()
}

func hashObject(filepath string) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	oid := HashObject(data)
	fmt.Println(oid)
}

func catFile(hashString string) {
	data := GetObject(hashString, "")
	_, err := os.Stdout.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

func writeTree() {
	fmt.Println(WriteTree())
}

func readTree(treeOid string) {
	ReadTree(treeOid)
}

func commit(message string) {
	fmt.Println(Commit(message))
}

func goLog() {
	oid, err := GetRef("HEAD")
	if err != nil {
		log.Fatal(err)
	}

	for oid != "" {
		commit := GetCommit(oid)

		fmt.Printf("commit %s\n", oid)
		fmt.Println(Indent(commit.message, "       "))

		oid = commit.parentOid
	}
}

func checkout(oid string) {
	Checkout(oid)
}

func tag(name string, oidArg string) {
	headOid, err := GetRef("HEAD")
	if err != nil {
		log.Fatal(err)
	}

	var oid string
	if oidArg == "" {
		oid = headOid
	}

	Tag(name, oid)
}
