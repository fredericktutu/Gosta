package main

/*
case4:

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
	
	go receive(ch)
	go receive(ch)
	
	ch <- 1  
	ch <- 1
	ch <- 1 // bug here 

	return 
}