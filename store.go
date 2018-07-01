package main

import (
	"github.com/mitchellh/go-homedir"
)

func getHome() string {
	h, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	d, err := homedir.Expand(h)
	if err != nil {
		panic(err)
	}
	return d + "/strmr/_data"
}
