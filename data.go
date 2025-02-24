// This is the main package
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// GitDir Git directory name
const GitDir string = ".gogit"

// Init Initializes the gogit directory
func Init() {
	os.Mkdir(GitDir, 0755)
	os.Mkdir(filepath.Join(GitDir, "objects"), 0755)
}

// HashObject hashes file contents to a unique blob
func HashObject(data []byte, typeArg ...string) (oid string) {
	dataType := "blob"
	if len(typeArg) > 0 {
		dataType = typeArg[0]
	}

	hash := sha256.New()
	hash.Write(data)
	byteHash := hash.Sum(nil)
	oid = hex.EncodeToString(byteHash)

	obj := []byte(dataType)
	obj = append(obj, 0)
	obj = append(obj, data...)

	os.WriteFile(filepath.Join(GitDir, "objects", oid), obj, 0444)
	return
}

// GetObject returns the contents of a file based on its hash
func GetObject(hash string, expectedTypeArg ...string) []byte {
	obj, err := os.ReadFile(filepath.Join(GitDir, "objects", hash))
	if err != nil {
		log.Fatal(err)
	}

	sepIndex := bytes.IndexByte(obj, 0)
	var dataType, content []byte

	dataType = obj[:sepIndex]
	content = obj[sepIndex+1:]

	if len(expectedTypeArg) > 0 {
		expectedDataType := expectedTypeArg[0]
		if string(dataType) != expectedDataType {
			log.Fatal(fmt.Errorf("Types are incompatible"))
		}
	}

	return content
}
