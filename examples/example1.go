package main

import (
	"fmt"
)

func main() {
	go A()
	go B()
	fmt.Println("baz")
}

func A() {
	fmt.Println("foo")
}

func B() {
	fmt.Println("bar")
}