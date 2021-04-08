package main

/*
[case3]
动态启动协程,但出错的方式和1 2一样

[expect] 18 Bugs
*/
import (
	"os"
	"strconv"
)

func send(ch chan int) {
	ch <- 1


}

func receive(ch chan int) {
	x, _ := strconv.ParseInt(os.Args[2], 10, 0)
	if x <= 1 {
		<- ch
	}

}

func main() {
	
	ch := make(chan int, 1)
	x, _ := strconv.ParseInt(os.Args[2], 10, 0)
	if x < 5 {
		go send(ch)
	}
	go receive(ch)
	
	y, _ := strconv.ParseInt(os.Args[3], 10, 0)
	if y < 2 {
		<- ch
	}
	return 
}