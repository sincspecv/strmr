package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

func getFileInfo(file map[string]string) {
	fmt.Println(file["dir"] + ": " + file["name"])
}

func getFileList(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	// Loop through each file to make sure it is a valid file or a
	// directory. If directory, recursively check new directory
	for _, f := range files {

		// Make sure directory has trailing slash
		lastChar := dir[len(dir)-1:]
		if lastChar != "/" {
			dir += "/"
		}

		// Check if we are dealing with a file or a directory, and act accordingly
		switch mode := f.Mode(); {

		case mode.IsDir():
			// If mode is a directory we need to check in this new directory
			fmt.Println("\t" + f.Name())
			getFileList(dir + f.Name())
			break

		case mode.IsRegular():
			// If it is a file, make sure it's an audio file and return data
			if filepath.Ext(dir+f.Name()) == ".mp3" || filepath.Ext(dir+f.Name()) == ".flac" || filepath.Ext(dir+f.Name()) == ".ogg" {
				file := make(map[string]string)
				file["dir"] = dir
				file["name"] = f.Name()
				getFileInfo(file)
			}

			break

		default:
			fmt.Println("not a valid file or directory")
			break
		}
	}
}

func scanDirectories() {
	directories := getDirectories()

	// Iterate through directories in config
	dirlen := len(directories)
	for i := 0; i < dirlen; i++ {
		fmt.Println(directories[i])
		getFileList(directories[i])
	}

}
