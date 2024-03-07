package main

import (
	"fmt"
	"log"
)

// find the logest leght of string
// find the leader from arr. leader means from its right all shound be less value
// last ele is always leader as there is no value. eg. a = [16,17,5,4,3] o/t:= 17,5,2
type Errorstr struct {
	HttpSts int
}

func (s *Errorstr) Error() string {
	return fmt.Sprintf("code %d", s.HttpSts)
}

type txtStr struct{}

func testStr() error {
	return &Errorstr{}
}

func update() error {
	return nil
}
func main() {
	fmt.Println("hello")
	err := testStr()
	if err != nil {
		log.Fatal("error from testStr", err)
	}
	if err = update(); err != nil {
		log.Fatal("error from update")
	}
	fmt.Println("success")
}
