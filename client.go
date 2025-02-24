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
				log.Fatal(fmt.Errorf("Invalid number of arguments"))
			}
			hashObject(args[0])
		},
	}

	var catFileCommand = &cobra.Command{
		Use:   "cat-file",
		Short: "Prints the content of an object",
		Run: func(_ *cobra.Command, args []string) {
			if len(args) > 1 {
				log.Fatal(fmt.Errorf("Invalid number of arguments"))
			}
			catFile(args[0])
		},
	}

	var writeTreeCommand = &cobra.Command{
		Use:   "write-tree",
		Short: "Writes the current working directory into object database",
		Run: func(_ *cobra.Command, args []string) {
			if len(args) > 1 {
				log.Fatal(fmt.Errorf("Invalid number of arguments"))
			}
			writeTree()
		},
	}

	var readTreeCommand = &cobra.Command{
		Use:   "read-tree",
		Short: "Reads the OID of a tree and writes the state to the current working directory",
		Run: func(_ *cobra.Command, args []string) {
			if len(args) > 1 {
				log.Fatal(fmt.Errorf("Invalid number of arguments"))
			}
			readTree(args[0])
		},
	}

	rootCmd.AddCommand(initCommand)
	rootCmd.AddCommand(hashObjectCommand)
	rootCmd.AddCommand(catFileCommand)
	rootCmd.AddCommand(writeTreeCommand)
	rootCmd.AddCommand(readTreeCommand)
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
	WriteTree()
}

func readTree(treeOid string) {
	ReadTree(treeOid)
}
