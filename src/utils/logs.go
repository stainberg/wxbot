package utils

import (
	"fmt"
	"time"
)

func Log(model, msg string) {
	fmt.Printf(`%s %s / %s`, time.Now().Format("2006-01-02 15:04:05"), model, msg)
	fmt.Println()
}
