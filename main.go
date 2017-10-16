package main

import (
	"builder/builder"
	"fmt"
)

func main() {
	builder.RunEnvironment(".", false)

	fmt.Println("press enter/return to quit")
	fmt.Scanln()
}
