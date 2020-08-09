package main

import (
	"time"
)

func main() {
	storage := &Storage{}
	storage.Initialize()

	storage.Set("k1", "v1", 10)
	storage.Set("k2", "v2", 2)
	storage.Set("k3", "v3", 5)

	time.Sleep(time.Duration(20) * time.Second)
}
