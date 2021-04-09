package main

/*
[case3]
哲学家进餐问题,奇数先左后右，偶数相反

[expect] 120
*/

func odd(left chan int, right chan int, lock chan int) {
	<- left
	<- right
	left <- 0
	right <- 0

	<- lock
}

func even(left chan int, right chan int, lock chan int) {
	<- right
	<- left
	left <- 0
	right <- 0

	<- lock
}


func main() {
	
	stick1 := make(chan int, 1)
	stick2 := make(chan int, 1)
	stick3 := make(chan int, 1)
	stick4 := make(chan int, 1)
	stick5 := make(chan int, 1)

	lock := make(chan int, 0)

	stick1 <- 1
	stick2 <- 2
	stick3 <- 3
	stick4 <- 4
	stick5 <- 5

	go odd(stick1, stick2, lock)  // 1 2
	go even(stick2, stick3, lock)  // 2 3
	go odd(stick3, stick4, lock)  // 3 4
	go even(stick4, stick5, lock)  // 4 5
	go odd(stick5, stick1, lock)  // 5 1

	lock <- 1
	lock <- 1
	lock <- 1
	lock <- 1
	lock <- 1

	return 
}