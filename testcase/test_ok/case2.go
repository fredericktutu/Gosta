package hello2
import "fmt"

// 加入无关控制流

func A(ch chan int) {
	ch <- 1
}

func use(a int) {
	fmt.Println(a)
}

func main() {
	ch := make(chan int, 1) 
	go A(ch)
	var x int
	x = <- ch
	x = x + 1
	if x == 1{
		use(x)
	}
}