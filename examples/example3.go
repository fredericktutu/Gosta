//协程死锁简单例子

package main

func main() {
	ch1 := make(chan int, 0)
	go A(ch1)
	<- ch1
	<- ch1
}

func A(ch1 chan int) {
	ch1 <- 1
}

