package main

/*
[case1]
两个goroutine之间的匹配
但是因为main直接结束了，其实检测不到

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
	go send(ch)
	go receive(ch)
	return 
}