package main

import (
	"github.com/go-co-op/gocron"
	"time"
)

func runCronJon() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(10).Seconds().Do(receiveData)

	s.StartAsync()
}
