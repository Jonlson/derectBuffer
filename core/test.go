package main

import (
	"fmt"
)

func Test() {
	s := []int{1, 2, 3, 4, 5, 6}
	fmt.Println(s)
	s1 := s[0:3]
	fmt.Println(s1)
	s2 := s[3:]
	fmt.Println(s2)
	s1 = append(s1, 1)
	fmt.Println(s, s1, s2)
}

func main() {
	Test()
}
