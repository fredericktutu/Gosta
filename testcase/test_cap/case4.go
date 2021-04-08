package main

/*
[case4]
2个有容量的channel

[expect] 2 Bugs
*/

func A(ch1 chan int, ch2 chan int) {
	<- ch2
	<- ch1  // 此处会解除 block1
	<- ch1  
	<- ch1  // 此处如果比main快就可以产生死锁


}

func main() {
	
	ch1 := make(chan int, 2)
	ch2 := make(chan int, 1)

	ch1 <- 1
	ch1 <- 1
	ch2 <- 1
	go A(ch1, ch2)
	ch1 <- 1  // 此处block1
	<- ch1  // 若此处没有A快，则会block2 -> 死锁
	return 
}