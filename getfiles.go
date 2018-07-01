package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/dhowden/tag"
	"github.com/eaburns/flac"
	mp3 "github.com/hajimehoshi/go-mp3"
)

// Media type defines all parameters of a particular type of media
type Media struct {
	format      string
	title       string
	artist      string
	album       string
	albumArtist string
	composer    string
	genre       string
	year        int
	disc        int
	discTotal   int
	track       int
	trackTotal  int
	length      int64
	location    string
}

func readBytes(r io.Reader, n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := io.ReadFull(r, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func get7BitChunkedInt(b []byte) int {
	var n int
	for _, x := range b {
		n = n << 7
		n |= int(x)
	}
	return n
}

func checkForFlac(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	offset := 47000
	b, err := readBytes(f, offset)
	if err != nil {
		panic(err)
	}

	size := get7BitChunkedInt(b[3:7])
	fmt.Println(size)

	if bytes.Contains(b, []byte("fLaC")) {
		return true
	} else {
		return false
	}
}

func getFileInfo(file map[string]string) {
	fullPath := file["dir"] + file["name"]
	f, err := os.Open(fullPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var _ io.Reader = (*os.File)(nil)
	// r, err := io.Reader(f)

	// use dhowden/tag to get the tags
	meta, err := tag.ReadFrom(f)
	if err != nil {
		panic(err)
	}
	fmt.Println(fullPath)
	fmt.Println(meta.Format())

	if meta.FileType() == "" {
		panic("Unknown format")
	}

	// TODO: get file format from file headers
	// Get the file format from file extension
	// e := filepath.Ext(fullPath)
	format := meta.FileType()

	// Get length of audio file
	var length int64

	if format == "MP3" {
		// first make sure it's not actually a FLAC file
		// with ID3 tags
		if checkForFlac(fullPath) {
			fm := flac.Decode(f)
			fmt.Println(fm)
		} else {
			// use go-mp3 to get length
			d, err := mp3.NewDecoder(f)
			if err != nil {
				panic(err)
			}
			length = d.Length()
		}
	}

	if format == "FLAC" {
		fm := flac.Decode(f)
		fmt.Println(fm)
	}

	disc, discTotal := meta.Disc()
	track, trackTotal := meta.Track()

	m := Media{
		format:      format,
		title:       meta.Title(),
		artist:      meta.Artist(),
		album:       meta.Album(),
		albumArtist: meta.AlbumArtist(),
		composer:    meta.Composer(),
		genre:       meta.Genre(),
		year:        meta.Year(),
		disc:        disc,
		discTotal:   discTotal,
		track:       track,
		trackTotal:  trackTotal,
		length:      length,
		location:    fullPath,
	}

	fmt.Println(m)
	// fmt.Println("\t" + meta.Title() + " By " + meta.Artist())
	// log.Print(meta.Picture())

	// myDBDir := getHome()
	// defer os.RemoveAll(myDBDir)
	// // (Create if not exist) open a database
	// store, err := db.OpenDB(myDBDir)
	// if err != nil {
	// 	panic(err)
	// }

	// // Create two collections: Feeds and Votes
	// if err := store.Create("Feeds"); err != nil {
	// 	panic(err)
	// }
	// if err := store.Create("Votes"); err != nil {
	// 	panic(err)
	// }

	// // What collections do I now have?
	// for _, name := range store.AllCols() {
	// 	fmt.Printf("I have a collection called %s\n", name)
	// }

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
			fmt.Println(f.Name())
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
