package main

/*
case4:

因为main中的send发生死锁，在发生死锁前有两种路径

[expect] 2 Bugs
*/


func send(ch chan int) {
	ch <- 1 
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