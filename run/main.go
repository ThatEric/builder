package main

import "builder"

func main() {
	readyChan := make(chan error)
	builder.RunEnvironment(readyChan)
}
