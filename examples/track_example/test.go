package main

import "fmt"

func main() {
Loop:
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			fmt.Println(i, j)
			if j == 2 {
				break Loop
			}
		}
	}
}
