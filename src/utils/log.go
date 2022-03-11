package utils

import (
	"fmt"
	"time"
)

func PrintError(err error) {
	fmt.Println("\n/!\\    Error    /!\\")
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("- - - - - - - - - -")
	fmt.Println(err)
}
