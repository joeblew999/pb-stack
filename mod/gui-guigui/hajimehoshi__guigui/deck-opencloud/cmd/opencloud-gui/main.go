package main

import (
	"fmt"
	"log"

	"opencloud/pkg/gui" // Assuming this path is correct based on original main.go
)

func main() {
	// No flags needed from the original set for gui.Start() currently
	// flag.Parse() // Not needed if no flags are defined for the GUI binary

	fmt.Println("OpenCloud - GUI Mode")
	fmt.Println("====================")
	fmt.Println("Starting GUI...")
	if err := gui.Start(); err != nil {
		log.Fatalf("GUI failed: %v", err)
	}
}
