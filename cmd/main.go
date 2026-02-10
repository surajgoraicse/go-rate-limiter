package main

import (
	"fmt"
	"os"
	"time"

	rate "github.com/surajgoraicse/go-rate-limiter/rate_v2"
)

func main() {
	limiter, _ := rate.New(10, 1)
	ip1 := "192.168.1.1"
	for range 20 {

		allow, err := limiter.Allow(ip1)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if allow {
			fmt.Println("success")
		} else {
			fmt.Println("rate limmited")
			time.Sleep(2 * time.Second)
		}
	}
	fmt.Println("compleet....")
}
