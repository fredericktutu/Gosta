package hello
import "fmt"
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
	use(x) 
	x = <- ch
	use(x)
}