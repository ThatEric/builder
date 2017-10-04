package main

import "github.com/thateric/builder"

func main() {
	readyChan := make(chan error)
	builder.RunEnvironment(readyChan)
}
