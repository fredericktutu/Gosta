package main

/*
[case4]
动态启动协程
Bug1:a b两处都正好成立，且最终选择了协程B和A的ch2匹配，这样A B正常结束，sync3.main卡死
Bug2:a 成立 但b不成立，这样B正常运行结束，sync2.A和sync3.main均卡死

[expect] 2 Bugs
*/



func A(ch1 chan int, ch2 chan int) {
	<- ch1  // sync 1
	var x int = 1
	if x >= 1 {  // a
		go B(ch1, ch2)
		<- ch1  //sync 2
	}
	<- ch2     // sync 3
	
}

func B(ch1 chan int, ch2 chan int) {
	var x int = 2
	if x <= 3 {  // b
		ch1 <- 2  //sync 2
		ch2 <- 3  //sync 3
	}

}

func main() {
	
	
	ch1 := make(chan int, 0)
	ch2 := make(chan int, 0)

	go A(ch1, ch2)
	ch1 <- 1  // sync 1
	ch2 <- 3  // sync 3
	

}