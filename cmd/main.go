package main

import (
	"fmt"
	"time"

	rate "github.com/surajgoraicse/token-bucket/rate_v1"
)

func main() {
	r := rate.New(10, 1)
	for range 20 {
		if r.Allow() {
			fmt.Println("success")
		} else {
			time.Sleep(2 * time.Second)
		}
	}
}
