package main

import (
	"fmt"
)

func factorial(n int) int {
	if n < 0 {
		return 0
	}
	result := 1
	if n > 1 {
		for i := 2; i <= n; i++ {
			result *= i
		}
	}
	return result
}

func main() {
	var n int
	fmt.Print("Enter a number to calculate its factorial: ")
	_, err := fmt.Scan(&n)
	if err != nil {
		fmt.Println("Invalid input. Please enter an integer.")
		return
	}

	result := factorial(n)
	fmt.Println("Factorial of", n, "is:", result)
}
