package main

import (
	"fmt"
)

func wait(n int) {
	for i := 1; i <= n*1000000; i++ {
		continue
	}
}

func A(ch1 chan int) {
	for i := 1; i < 5; i++ {
		fmt.Println("A", i)
		wait(1)
	}
	ch1 <- 1

}

func B(ch1, ch2 chan int) {
	fmt.Println("B gets", <-ch1)
	for i := 1; i < 5; i++ {
		fmt.Println("B", i)
		wait(1)
	}
	ch2 <- 2
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go A(ch1)      //将函数A制作为协程并启动
	go B(ch1, ch2) //将函数B制作为协程并启动

	fmt.Println("Main gets", <-ch2)

}
