package main

import (
	"fmt"
)


func A(lst *[]int) {

	for i, _ := range *lst {
		(*lst)[i] = 3
	}
}

func main() {
	var lst []int
	lst = append(lst, 1)
	lst = append(lst, 2)

	A(&lst)
	fmt.Println(lst)
}