package main
import (
	"fmt"
)
type Person struct {
	Name string
}

type Student struct {
	Person
	SId int
}

func main() {
	var stu Student
	var ps *Student
	ps = &stu
	stu.Name = "名字"
	stu.SId = 11
	fmt.Println(ps.Name) 


}