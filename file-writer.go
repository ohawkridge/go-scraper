package main

import (
	"os"
	"time"
)

func writeTxt() {
	ts := time.Now().UTC().String()
	err := os.WriteFile("tmp/test.txt", []byte("Last run:\n"+ts), 0755)
	if err != nil {
		panic(err)
	}
}
