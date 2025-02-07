package main

import (
	"fmt"
	"log"
)

var (
	buildDate   string
	buildCommit string
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	fmt.Println("MyHealthBot")
	return nil
}
