package main

/*
case2:
因为send产生死锁

expect: 1 Bugs
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