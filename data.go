// This is the main package
package main

import (
	"os"
)

// GitDir Git directory name
const GitDir string = ".gogit"

// Init Initializes the gogit directory
func Init() {
	os.Mkdir(GitDir, 0755)
}
