package main

/*
[case2]
容量为1，main和五个协程同时竞争，因此有5种错法

[expect] 5 Bugs
*/

func send(ch chan int) {
	ch <- 1
}

func receive(ch chan int) {
	<- ch
}

func main() {
	
	ch := make(chan int, 1)
	go send(ch)
	go send(ch)
	go send(ch)
	go send(ch)
	go send(ch)
	ch <- 1   //若这里是最后一个安排到的，则出错
	return 
}