package strmr

import (
	"fmt"
)

func main() {
	directories := getDirectories()

	// Iterate through directories in config
	dirlen := len(directories)
	for i := 0; i < dirlen; i++ {
		fmt.Println(directories[i])
	}
}
