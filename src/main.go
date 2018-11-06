package main

import "fmt"
import "os"
import "./repl"

func main() {

	panic("TEST")

	fmt.Printf("Hello! This is the Monkey programming language!\n")
	fmt.Printf("Feel free to type commands\n")

	repl.Start(os.Stdin, os.Stdout)
}
