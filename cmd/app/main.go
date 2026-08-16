package main

import "simpletorrent/internal/utils"
import "fmt"

func main() {
	// byteArray := []byte("li32e6:codinge")
	byteArray := []byte("d1:ai1e1:b3:fooe")
	str, nextInd, err := utils.ParseDict(byteArray)
	str.PrintValue()
	fmt.Println("NextIndex: ", nextInd)
	fmt.Println(err)
}