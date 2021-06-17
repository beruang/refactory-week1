package main

import (
	"github.com/prometheus/common/log"
	"refactory/notes/cmd"
)

func main() {
	if err := cmd.Start(); nil != err {
		log.Fatal(err)
	}
}
