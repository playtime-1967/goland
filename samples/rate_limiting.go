package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

// $@#!%*golang
func RateLimiter() {

	limiter := rate.NewLimiter(2, 3) // 2 events/sec, burst of 3
	for i := 1; i <= 10; i++ {
		err := limiter.Wait(context.TODO()) // blocks until permitted
		if err != nil {
			panic(err)
		}
		fmt.Printf("Request %d allowed at %v\n", i, time.Now())
	}
}
