package main

import (
	"log"
	"os"

	"github.com/NikolaB131/computer-club-system/internal/app"
)

type Flags struct {
	InputFilePath string
}

func main() {
	flags := mustParseFlags()

	appInstance := app.App{InputFilePath: flags.InputFilePath}
	appInstance.Run()
}

func mustParseFlags() Flags {
	if len(os.Args) < 2 {
		log.Fatal("First argument must be specified with the path to input file")
	}
	inputFilePath := os.Args[1]

	return Flags{InputFilePath: inputFilePath}
}
