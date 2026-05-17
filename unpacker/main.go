package main

import (
	"fmt"
	"unpacker/unpack"
)

func main() {
	e, err := unpack.Unpack(`mar5at`)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(e)
}
