package main

import (
	// Assuming your module path is 'github.com/joeblew999/pb-stack'
	// and main.go is in 'mod/reck/', the cmd package is in 'mod/reck/cmd/'.
	"github.com/joeblew999/pb-stack/mod/reck/cmd"
)

func main() {
	// Call the Execute function from the cmd package
	cmd.Execute()
}
