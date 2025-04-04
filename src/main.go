package main

import (
	"fmt"
	"os"
	"strings"
	"calcWithTests/src/calculator"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("no argument provided")
		fmt.Println("usage: ./calculate [your expression]")
		return
	}

	if strings.Compare(args[0], "-h") == 0 || strings.Compare(args[0], "--help") == 0 {
		fmt.Println("usage: ./calculate [your expression]")
		return
	}
	if len(args) != 1 {
		fmt.Println("no argument provided")
		fmt.Println("usage: ./calculate [your expression]")
		return
	}
	expr := args[0]
	result, err := calculator.Calculate(expr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(result)
}
