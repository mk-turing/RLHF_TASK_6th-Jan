package main

import "fmt"

func factors(n int) ([]int, error) { factors := []int{} if n < 2 { return nil, fmt.Errorf("Invalid input: %d must be greater than or equal to 2", n) } for i := 1; i <= n; i++ { if n%i == 0 { factors = append(factors, i) } } return factors, nil }

func main() { n := 1000000 _, err := factors(n) if err != nil { fmt.Println(err) } }