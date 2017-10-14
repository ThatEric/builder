package main

import "builder/builder"

func main() {
	readyChan := make(chan error)
	builder.RunEnvironment(".", readyChan)
}
