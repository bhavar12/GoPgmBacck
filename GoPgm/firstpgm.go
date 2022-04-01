package main

import "fmt"

type T struct {
	a int
}

func (tv T) Mv(a int) int { return 0 }

func getSun() {
	return ""
}
func main() {
	var t T
	fmt.Println(t.Mv(7))
	fmt.Println(T.Mv(t, 7))
	fmt.Println((T).Mv(t, 7))
	f1 := T.Mv
	fmt.Println(f1(t, 7))
	f2 := (T).Mv
	fmt.Println(f2(t, 7))
}
