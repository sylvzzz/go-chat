package main

import ("fmt"
		"time"
		"strconv")

func main() {
	var integers int = 1
	var floats float64 = 67.69
	var boolean bool = true
	var strings string = "The python with C speed?"
	var months []string = []string{"October", "December", "January", "April", "September"}  // literal
	// var names_alloced = make([]string, 0, 4) this is an array ready with allocated space
	student := make(map[string]any)

	student["name"] = "dbotelho"
	student["rank"] = "4"
	student["age"] = "20"

	fmt.Println("=== Standard types ===\n")
	fmt.Printf(" integers: %v \n floats: %v \n boolean: %v \n strings: %q\n", integers, floats, boolean, strings)
	fmt.Println("\nHello world!")
	fmt.Println("The time is", time.Now())

	age, err := strconv.Atoi(student["age"].(string))

	if err != nil {
		fmt.Printf("Invalid age of %v ...", student["age"])
	}

	if age < 21 {
		fmt.Printf("\n%v is one of the youngest students in 42 Lisbon with %v years of age, he´s also Rank %v\n", student["name"], student["age"], student["rank"])
	}

	for i := 0; i < 5; i++ {
		fmt.Printf("%v was rank %v at 42 Lisbon in %v\n",student["name"], i, months[i])
	}
}
