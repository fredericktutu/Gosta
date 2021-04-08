package main

/*
[case1]
所有的x都是不确定的，所有的if路径组合都可能发生

-> 在这里把权值从1000调到10000了

[expect] 82 Bugs
*/
import (
	"os"
	"strconv"
)

func send(ch chan int) {
	x, _ := strconv.ParseInt(os.Args[1], 10, 0)
	if x <= 1 {
		ch <- 1
	}

}

func receive(ch chan int) {
	x, _ := strconv.ParseInt(os.Args[2], 10, 0)
	if x <= 1 {
		<- ch
	}

}

func main() {
	
	ch := make(chan int, 1)
	go send(ch)
	go receive(ch)
	
	x, _ := strconv.ParseInt(os.Args[3], 10, 0)
	if x < 2 {
		<- ch
	}
	return 
}