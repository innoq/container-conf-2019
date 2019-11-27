package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	size, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	arr := make([]byte, size*1024*1024)
	fmt.Printf("Allocated: %T, size: %v megabyte\n", arr, size)
}
