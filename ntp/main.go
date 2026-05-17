package main

import (
	"fmt"
	"os"

	"github.com/beevik/ntp"
)

func main() {
	currentTime, err := ntp.Time("time.google.com")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(currentTime.Local())
}
