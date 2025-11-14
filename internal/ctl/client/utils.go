package client

import (
	"fmt"
	"time"
)

func Log(msg string) {
	fmt.Printf("[%s] %s\n", time.Now().Format("15:04:05.000"), msg)
}
