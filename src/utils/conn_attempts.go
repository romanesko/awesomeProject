package utils

import (
	"log"
	"time"
)

func ConnectionAttemps(conn_func func() error, attemps int, delay time.Duration) (err error) {
	for i := 0; i < attemps; i++ {
		err = conn_func()
		if err != nil {
			log.Printf("attempting to connect: current attemp: %d, appemps left: %d", i+1, attemps-i-1)
			time.Sleep(delay)
			continue
		}
	}
	return err
}
