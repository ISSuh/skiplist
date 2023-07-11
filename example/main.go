package main

import (
	"fmt"
)

func main() {
	fmt.Println("skiplist example")

	list := skipList.New(5)
	list.Set("test", "test")
	list.Print()

	item := list.Get("test")
	println("item: ", item)
}
