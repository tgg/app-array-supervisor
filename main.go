package main

import "fmt"

const (
	Unknown = iota
	Starting
	Started
	Running
	Stopping
	Stopped
)

func main() {
	fmt.Println("app-array-supervisor")
}
