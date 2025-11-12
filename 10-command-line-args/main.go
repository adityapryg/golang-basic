package main

import (
	"flag"
	"fmt"
	"os"
)

// calculator function for basic math operations
func calculator(a float64, b float64, operation string) float64 {
	switch operation {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	case "/":
		if b != 0 {
			return a / b
		}
		fmt.Println("Error: Division by zero")
		return 0
	default:
		fmt.Println("Error: Unsupported operation")
		return 0
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  go run main.go osargs | flags | calc | server")
		return
	}

	switch os.Args[1] {
	case "osargs":
		fmt.Println("All arguments:", os.Args)
	case "flags":
		var name string
		var age int
		var verbose bool
		flag.StringVar(&name, "name", "", "Your Name")
		flag.IntVar(&age, "age", 0, "Your Age")
		flag.BoolVar(&verbose, "verbose", false, "Enable verbose mode")
		flag.Parse()
		if verbose {
			fmt.Printf("Verbose Mode: Name: %s, Age: %d\n", name, age)
		} else {
			fmt.Printf("Name: %s, Age: %d\n", name, age)
		}
	case "calc":
		if len(os.Args) < 5 {
			fmt.Println("Usage: go run main.go calc <num1> <num2> <operation>")
			return
		}
		num1 := os.Args[2]
		num2 := os.Args[3]
		operation := os.Args[4]
		result := calculator(num1, num2, operation)
		fmt.Printf("Result: %f\n", result)
	case "server":
		var host string
		var port int
		var debug bool
		flag.StringVar(&host, "host", "localhost", "Server Host")
		flag.IntVar(&port, "port", 8080, "Server Port")
		flag.BoolVar(&debug, "debug", false, "Enable Debug Mode")
		flag.Parse()
		fmt.Printf("Starting server on %s:%d with debug=%t\n", host, port, debug)
	default:
		fmt.Println("Error: Unknown command")
	}
}