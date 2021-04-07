package main

/*
case1:
正常存取
expect: 0 Bugs
*/

func send(ch chan int) {
	ch <- 1
}

func receive(ch chan int) {
	<- ch
}

func main() {
	
	ch := make(chan int, 0)
	go send(ch)
	go receive(ch)
	return 
}