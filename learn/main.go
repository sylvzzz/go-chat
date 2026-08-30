package main

import ("fmt"
		"time")

func main() {
	var c, python = "c", "python"
	var integers int = 1
	var floats float64 = 67.69
	var boolean bool = true
	var strings string = "The python with C speed?"

	fmt.Println("=== Standard types ===\n")
	fmt.Printf(" integers: %v \n floats: %v \n boolean: %v \n strings: %q\n", integers, floats, boolean, strings)

	fmt.Println("\nHello world!")
	fmt.Println("The time is", time.Now())
}
