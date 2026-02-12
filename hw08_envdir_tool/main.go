package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
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

	args = args[1:]
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = []string{}
	var out strings.Builder
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out.String())
}
