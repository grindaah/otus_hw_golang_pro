package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

const helloOtus = "Hello, OTUS!"

func main() {
	fmt.Printf("%s", stringutil.Reverse(helloOtus))
}
