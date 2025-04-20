package main

import (
	"calcWithTests/src/calculator"
	"flag"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("no argument provided")
		fmt.Println("usage: ./calculate [your expression]")
		return
	}

	helpFlag := flag.Bool("help", false, "usage: ./calculate [your expression]")
	angleUnit := flag.String("angle-unit", "radian", "Angle unit (degree or radian)")

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if *angleUnit != "degree" && *angleUnit != "radian" {
		fmt.Println("Error: angle-unit must be either 'degree' or 'radian'")
		flag.Usage()
		os.Exit(1)
	}

	expr := args[len(args)-1]
	result, err := calculator.Calculate(expr, calculator.CalculatorConfig{AngleUnits: *angleUnit})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(result)
}
