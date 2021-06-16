package main

import (
	"github.com/prometheus/common/log"
	"refactory/notes/internal/server"
)

func main() {
	if err := server.Start(); nil != err {
		log.Fatal(err)
	}
}
