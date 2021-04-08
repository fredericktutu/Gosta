package main

/*
[case1]
容量为1
S->R->S 不产生bug

[expect] 0 Bugs
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
	go receive(ch)
	
	ch <- 1
	return 
}