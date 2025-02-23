package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type entry struct {
	dataType string
	oid      string
	name     string
}

// WriteTree Write the current working directory into the object store
func WriteTree(targetDir ...string) (oid string) {
	directory := "."
	if len(targetDir) > 0 {
		directory = targetDir[0]
	}
	entries := make([]entry, 1)

	filepath.WalkDir(directory, func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		var dataType string
		var oid string
		if !isIgnored(s) {
			if !d.IsDir() {
				dataType = "blob"
				data, err := os.ReadFile(s)
				if err != nil {
					log.Fatal(err)
				}
				oid = HashObject(data)
			} else {
				dataType = "tree"
				oid = WriteTree(s)
			}

			entry := entry{dataType: dataType, oid: oid, name: s}
			entries = append(entries, entry)
		}
		return nil
	})

	formattedEntries := make([]string, 0, len(entries))
	for _, n := range entries {
		formattedEntries = append(formattedEntries, fmt.Sprintf("%s %s %s", n.dataType, n.oid, n.name))
	}
	tree := strings.Join(formattedEntries, "")
	return HashObject([]byte(tree), "tree")
}

func isIgnored(path string) bool {
	files := strings.SplitSeq(path, "/")
	for f := range files {
		if f == ".gogit" {
			return true
		}
	}
	return false
}
