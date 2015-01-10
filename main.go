/*
Safe rm

Run as saferm <filename>

Moves file to .safetrash in home directory.
*/
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println(os.Args[1])
}
