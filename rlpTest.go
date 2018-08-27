package main

import (
	"fmt"

	"github.com/go-ShardingBlockchain/rlp"
)

type Person struct {
	ID       string
	Position string
}

type User struct {
	Name string
	Age  uint
	Test uint
	P    Person
}

func main() {
	user := &User{
		Name: "bbbbbbbb",
		Age:  20,
		Test: 19,
		P: Person{
			ID:       "123456",
			Position: "shanghai shi",
		},
	}

	bt, _ := rlp.EncodeToBytes(user)
	fmt.Printf("%v->%x\n", user, bt)

	var user1 = new(User)
	rlp.DecodeBytes(bt, &user1)
	fmt.Printf("%x -> %v\n", bt, user1)
}
