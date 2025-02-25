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

// CommitDetails Holds details of a commit
type CommitDetails struct {
	treeOid   string
	parentOid string
	message   string
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

// Commit commits the current state
func Commit(message string) string {
	var commitString string = fmt.Sprintf("tree %s\n", WriteTree())

	HEAD, err := GetRef("HEAD")
	if err != nil {
		log.Fatal(err)
	}

	if HEAD != "" {
		commitString += fmt.Sprintf("parent %s\n", HEAD)
	}

	commitString += "\n"
	commitString += fmt.Sprintf("%s\n", message)

	commitOid := HashObject([]byte(commitString), "commit")
	UpdateRef("HEAD", commitOid)
	return commitOid
}

// GetCommit Returns the commit details of a given commit OID
func GetCommit(oid string) CommitDetails {
	var parent, tree, message string
	commit := string(GetObject(oid, "commit"))
	commitLines := strings.Split(commit, "\n")
	for _, line := range commitLines {
		if line == "" {
			break
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			log.Fatal("Not enough parts to commit detail")
		}

		key, val := parts[0], parts[1]
		switch key {
		case "tree":
			tree = val
		case "parent":
			parent = val
		default:
			log.Fatal(fmt.Sprintf("Unknown key in commit details: %s", key))
		}
	}

	message = strings.Join(commitLines[len(commitLines)-2:], "\n")

	return CommitDetails{treeOid: tree, parentOid: parent, message: message}
}

// Checkout Check out a specified commit
func Checkout(oid string) {
	commit := GetCommit(oid)
	ReadTree(commit.treeOid)
	UpdateRef("HEAD", oid)
}

// Tag Tags a specified oid with a given name
func Tag(name string, oid string) {
	UpdateRef(filepath.Join("refs", "tags", name), oid)
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
