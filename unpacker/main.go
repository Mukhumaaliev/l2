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

// func unpack(input string) (string, error) {
// 	if "" == input {
// 		return "", nil
// 	}
// 	str := []rune(input)
// 	isdigit := false
// 	var result []rune
// 	var prev rune
// 	var count int
// 	count_cifer := 0

// 	for _, v := range str {
// 		if v >= '0' && v <= '9' {
// 			count_cifer++
// 		}

// 		if v < '0' || v > '9' {
// 			isdigit = true
// 		} else {
// 			isdigit = false
// 			break
// 		}
// 	}
// 	if !isdigit {
// 		return input, nil
// 	}
// 	if count_cifer == len(str) {
// 		return "", fmt.Errorf("invalid string")
// 	}

// 	for ind, v := range str {

// 		if ind == 0 {
// 			prev = v
// 		}
// 		if v >= '0' && v <= '9' && prev == '\'' {
// 			prev = v
// 		} else if v >= '0' && v <= '9' && prev != '\'' {
// 			count = int(v)
// 			for i := 0; i < count; i++ {
// 				result = append(result, prev)
// 			}
// 		} else if v != '\'' && v < '0' || v > '9' && str[ind+1] >= '0' && str[ind+1] <= '9' {
// 			count = int(str[ind+1])
// 			for i := 0; i < count; i++ {
// 				result = append(result, prev)
// 			}
// 		} else {
// 			if v != '\'' {
// 				result = append(result, v)
// 			}
// 			prev = v
// 		}
// 	}
// 	return string(result), nil
// }
