package main

/*
[case1]
sender & receiver

[expect] 137
*/

func sender(ch chan int) {
	ch <- 1
	ch <- 1
	ch <- 1
	ch <- 1
	ch <- 1
	ch <- 1
	ch <- 1
}


func receiver1(ch chan int, ch2 chan int) {
	<- ch
	<- ch
	<- ch
	<- ch
	<- ch2
}
func receiver2(ch chan int, ch2 chan int) {
	<- ch
	<- ch
	<- ch
	<- ch2
}


func main() {
	
	ch := make(chan int, 7)
	ch2 := make(chan int, 0)
	go sender(ch)
	var x = 1
	receiver1(ch, ch2)
	if x > 0 {
		go receiver1(ch, ch2)
	} else {
		go receiver2(ch, ch2)
	}
	
	ch2 <- 0
	ch2 <- 0
	
	return 
}