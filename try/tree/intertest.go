package intertest

import (
	"fmt"
)
type Animal interface {
	Eat();
}

type Cat struct {
	Name string;
	Catid int;
}

func (cat Cat) Eat() {
	fmt.Println(cat.Name, cat.Catid, "eating")
}

type Dog struct {
	Name string;
	Dogid int;
}

func (dog Dog) Eat() {
	fmt.Println(dog.Name, dog.Dogid, "eating")
}

func PrintEat(animal Animal) {
	switch animal.(type) {
	case Cat:
		fmt.Println("a cat")
		animal.Eat()
	case Dog:
		fmt.Println("a dog")
		animal.Eat()
	}
}

func main() {
	cat := Cat{
		Name: "cat1",
		Catid: 1,
	}
	dog := Dog{
		Name: "dog1",
		Dogid: 1,
	}
	PrintEat(cat)
	PrintEat(dog)

}