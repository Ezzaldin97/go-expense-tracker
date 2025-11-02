package src

import (
	"os"
)

func CreateDir(name string) {
	err := os.MkdirAll(name, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
