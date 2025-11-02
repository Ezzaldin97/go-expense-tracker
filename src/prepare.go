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

func remove[T comparable](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}
