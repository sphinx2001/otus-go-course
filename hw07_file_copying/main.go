package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	// Определить длину файла
	fileInfo, err := os.Stat(from)
	if err != nil {
		fmt.Println("Ошибка получения информации о файле:", err)
		return
	}
	// определить задавались ли параметры offset и limit
	provided := make(map[string]bool)

	flag.Visit(func(f *flag.Flag) {
		provided[f.Name] = true
	})

	if !provided["limit"] {
		limit = fileInfo.Size()
	}
	if !provided["offset"] {
		offset = 0
	}

	err = Copy(from, to, offset, limit)
	if err != nil {
		fmt.Println("Ошибка копирования:", err)
	}
}
