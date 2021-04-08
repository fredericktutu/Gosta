package main

/*
[case2]
因为send而产生死锁，有两种可能

[expect] 2 Bugs
*/


func send(ch chan int) {
	ch <- 1  //在这里block住没有关系
}

func receive(ch chan int) {
	<- ch
}

func main() {
	
	ch := make(chan int, 0)
	

	go send(ch)
	go receive(ch)

	ch <- 1  // 这里block住会有死锁

	return 
}