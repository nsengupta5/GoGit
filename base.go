package main

import (
	"fmt"
	"log"
	"maps"
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

// ReadTree takes the OID of a tree and extracts the contents to the working directory
func ReadTree(treeOid string, basePathArg ...string) {
	emptyCurrentDirectory()
	basePath := "./"
	if len(basePathArg) > 0 {
		basePath = basePathArg[0]
	}
	treeMap := getTree(treeOid, basePath)
	for path, oid := range treeMap {
		os.MkdirAll(filepath.Dir(path), 0755)
		os.WriteFile(path, GetObject(oid), 0644)
	}
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

func iterTreeEntries(oid string) ([]entry, error) {
	if oid == "" {
		return nil, nil
	}

	tree := string(GetObject(oid, "tree"))
	treeEntry := strings.Split(tree, "\n")

	var entries []entry
	for _, entryLine := range treeEntry {
		if entryLine == "" {
			continue
		}

		parts := strings.SplitN(entryLine, " ", 3)
		if len(parts) < 3 {
			log.Fatal(fmt.Errorf("Not enough data for given object in the tree"))
		}

		entryObj := entry{dataType: parts[0], oid: parts[1], name: parts[2]}
		entries = append(entries, entryObj)
	}

	return entries, nil
}

func getTree(oid string, basePathArg ...string) map[string]string {
	result := make(map[string]string)
	basePath := ""
	if len(basePathArg) > 0 {
		basePath = basePathArg[0]
	}

	treeEntries, _ := iterTreeEntries(oid)
	for _, entry := range treeEntries {
		if strings.Contains(entry.name, "/") {
			log.Fatal("Invalid tree: path contains /")
		}
		if strings.Contains(entry.name, "./") || strings.Contains(entry.name, "../") {
			log.Fatal("Invalid tree: path contains . or ..")
		}
		path := basePath + entry.name
		if entry.dataType == "blob" {
			result[path] = entry.oid
		} else if entry.dataType == "tree" {
			subtree := getTree(oid, fmt.Sprintf("%s/", path))
			maps.Copy(result, subtree)
		} else {
			log.Fatal(fmt.Sprintf("Unknown tree entry %s", entry.dataType))
		}
	}
	return result
}

func emptyCurrentDirectory() {
	dirEntries, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, de := range dirEntries {
		path := filepath.Join(".", de.Name())
		if isIgnored(path) {
			continue
		}

		err = os.RemoveAll(de.Name())
		if err != nil {
			log.Fatal(err)
		}
	}
}
