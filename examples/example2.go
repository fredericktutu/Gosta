package main

import (
	"fmt"
)

func main() {
	ch1 := make(chan int, 0)
	ch2 := make(chan int, 0)
	go A(ch1)
	go B(ch1, ch2)
	<- ch2
	fmt.Println("baz")


}

func A(ch1 chan int) {
	fmt.Println("foo")
	ch1 <- 1
}

func B(ch1 chan int, ch2 chan int) {
	<- ch1
	fmt.Println("bar")
	ch2 <- 1
}