package main

/*
case3:

两次S和R都能匹配上

[expect] 0 Bugs
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
	



	return 
}