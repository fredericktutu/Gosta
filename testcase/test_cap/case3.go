package main

/*
[case3]
容量为5，main和五个协程同时竞争，因此错法数为5的全排列，只要main在最后即可

[expect] 5 Bugs
*/

func send(ch chan int) {
	ch <- 1
}

func receive(ch chan int) {
	<- ch
}

func main() {
	
	ch := make(chan int, 5)
	go send(ch)
	go send(ch)
	go send(ch)
	go send(ch)
	go send(ch)
	ch <- 1   //若这里是最后一个安排到的，则出错
	return 
}