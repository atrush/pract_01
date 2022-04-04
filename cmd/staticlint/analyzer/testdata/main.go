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

	if qq == 43 {
		os.Exit(0) // want "os.Exit call in main.go"
	}

	fmt.Printf("n%v", qq)

	//exit
	os.Exit(0) // want "os.Exit call in main.go"
}
