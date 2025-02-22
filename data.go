// This is the main package
package main

import (
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
func HashObject(data []byte) {
	hash := sha256.New()
	hash.Write(data)
	byteHash := hash.Sum(nil)
	stringHash := hex.EncodeToString(byteHash)

	os.WriteFile(filepath.Join(GitDir, "objects", stringHash), data, 0444)
	fmt.Println(stringHash)
}

// GetObject returns the contents of a file based on its hash
func GetObject(hash string) (data []byte) {
	data, err := os.ReadFile(filepath.Join(GitDir, "objects", hash))
	if err != nil {
		log.Fatal(err)
	}
	return
}
