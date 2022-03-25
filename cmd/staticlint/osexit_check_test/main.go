package main

import (
	"fmt"
	"os"
)

func main() {

	qq := 9
	for i := 0; i < 9; i++ {
		qq = 2
	}
	fmt.Printf("n%v", qq)
	os.Exit(0)

	//exit
	os.Exit(0)
}
