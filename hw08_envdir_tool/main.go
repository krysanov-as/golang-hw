package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Println("Usage example: go-envdir /path/to/env/dir command arg1 arg2")
		os.Exit(baseMainError)
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Println("Error reading envdir:", err)
		os.Exit(baseMainError)
	}

	returnCode := RunCmd(os.Args[2:], env)
	os.Exit(returnCode)
}
