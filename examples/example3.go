package main

import (
	"fmt"
)

func wait(n int) {
	for i := 1; i <= n*1000000; i++ {
		continue
	}
}

func A(ch chan int) {
	for i := 1; i < 5; i++ {
		fmt.Println("A", i)
		wait(1)
	}
	ch <- 1

}

func B(ch chan int) {
	for i := 1; i < 5; i++ {
		fmt.Println("B", i)
		wait(1)
	}
	ch <- 2
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go A(ch1)
	go B(ch2)
	var x int
	select {
	case x = <-ch1:
		fmt.Println("A is selected, x is ", x)
	case x = <-ch2:
		fmt.Println("B is selected, x is ", x)
	default:
		fmt.Println("stuck!")
	}

}
