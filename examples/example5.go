package main

func main() {
	var x

	if x >= 10 {
		go A()
	}
	if x < 5 {
		go B()
	}
}