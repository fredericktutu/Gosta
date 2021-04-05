// 协程自动机的例子

package main

func main() {
	var x int = 1

	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)

	go A(ch1)

	if x <= 1 {
		go A(ch2)
	}

	ch1 <- 1
	ch2 <- 2
}

func A(ch chan int) {
	<- ch
}