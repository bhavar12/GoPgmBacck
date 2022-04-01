package main

import "fmt"

type Guitarst interface {
	PlayGuitar()
}

type BaseGuitar struct {
	Name string
}

type VeryGuitar struct {
	Name string
}

func (a BaseGuitar) PlayGuitar() {
	fmt.Println(a.Name)
}
func (b VeryGuitar) PlayGuitar() {
	fmt.Println(b.Name)
}

func add(a interface{}) {
	fmt.Println(a)
}

func main() {
	test := 10
	add(test)
	str := "test something"
	add(str)

	var player BaseGuitar
	player.Name = "Hello"
	player.PlayGuitar()

	var player1 VeryGuitar

	player1.Name = "World"
	player1.PlayGuitar()

}
