package src

import (
	"os"
)

func LoggerDir() {
	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		panic(err)
	}
}
