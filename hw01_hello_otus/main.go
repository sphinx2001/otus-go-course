package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	s := "Hello, OTUS!"
	s = reverse.String(s)
	fmt.Println(s)
}
