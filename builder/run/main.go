package main

import "github.com/cakemarketing/CapService/builder"

func main() {
	readyChan := make(chan error)
	builder.RunEnvironment(readyChan)
}
