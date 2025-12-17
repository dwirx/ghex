package main

import (
	"fmt"
	"os"
)

// Legacy ghe command - redirects to ghex
func main() {
	fmt.Println("ghe has been renamed to ghex")
	fmt.Println("Please use 'ghex' instead")
	os.Exit(0)
}
