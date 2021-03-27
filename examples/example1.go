package main

import (
	"fmt"
)

func wait(n int) {
	for i := 1; i <= n*1000000; i++ {
		continue
	}
}

func A() {
	for i := 1; i < 5; i++ {
		fmt.Println("A", i)
		wait(1)
	}

}

func B() {
	for i := 1; i < 5; i++ {
		fmt.Println("B", i)
		wait(1)
	}
}

func main() {
	go A() //将函数A制作为协程并启动
	go B() //将函数B制作为协程并启动
	wait(10)

}
