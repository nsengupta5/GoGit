package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
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
	entries := make([]entry, 0)

	dirEntries, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	for _, de := range dirEntries {

		path := filepath.Join(directory, de.Name())
		if isIgnored(path) {
			continue
		}

		var dataType string
		var oid string

		if !de.IsDir() {
			dataType = "blob"
			data, err := os.ReadFile(path)
			if err != nil {
				log.Fatal(err)
			}
			oid = HashObject(data)
		} else {
			dataType = "tree"
			oid = WriteTree(path)
		}

		entry := entry{dataType: dataType, oid: oid, name: de.Name()}
		entries = append(entries, entry)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].name < entries[j].name
	})

	formattedEntries := make([]string, 0, len(entries))
	for _, n := range entries {
		formattedEntries = append(formattedEntries, fmt.Sprintf("%s %s %s", n.dataType, n.oid, n.name))
	}
	tree := strings.Join(formattedEntries, "\n")
	return HashObject([]byte(tree), "tree")
}

func isIgnored(path string) bool {
	files := strings.SplitSeq(path, "/")
	for f := range files {
		if f == ".gogit" || f == ".git" {
			return true
		}
	}
	return false
}
