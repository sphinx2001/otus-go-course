package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Empty args...")
		return
	}
	envs, err := ReadDir(args[0])
	if err != nil {
		fmt.Printf("ReadDir error: %v\n", err.Error()+"\n")
		return
	}

	os.Exit(RunCmd(args[1:], envs))
}
