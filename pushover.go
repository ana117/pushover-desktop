package main

import (
	"fmt"

	"github.com/gen2brain/beeep"
)

func main() {
	fmt.Println("Hello, World!")
	beeep.AppName = "My App Name"

	err := beeep.Alert("Title", "Message body", "testdata/warning.png")
	if err != nil {
		panic(err)
	}
}
