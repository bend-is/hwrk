package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}

	envs, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Printf("Failed to read environments from directory: %s\n", err)
		os.Exit(1)
	}

	os.Exit(RunCmd(os.Args[2:], envs))
}
