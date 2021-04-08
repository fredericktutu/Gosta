package main

/*
[case2]
和case1相同的控制流，但是如果加上SMT，有些路径就不可能
预期结果应该和case1一样

[expect] 82 Bugs
*/

func send(ch chan int) {
	x := 1
	if x <= 1 {
		ch <- 1
	}

}

func receive(ch chan int) {
	x := 1
	if x <= 1 {
		<- ch
	}

}

func main() {
	
	ch := make(chan int, 1)
	go send(ch)
	go receive(ch)
	
	x := 1
	if x < 2 {
		<- ch
	}
	return 
}