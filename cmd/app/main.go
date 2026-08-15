package main

import "simpletorrent/internal/utils"
import "fmt"

func main() {
	byteArray := []byte("i3200000000000000000000000000000000000000e")
	str, nextInd, err := utils.ParseInteger(byteArray)
	str.PrintValue()
	fmt.Println("NextIndex: ", nextInd)
	fmt.Println(err)
}