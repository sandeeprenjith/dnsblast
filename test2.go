package main

import "time"
import "os"

func main() {
	var n int
	n = 10
	ender := time.Tick(time.Duration(n) * time.Second)

	go func() {
		<-ender
		os.Exit(0)
	}()
	time.Sleep(60 * time.Second)
}
