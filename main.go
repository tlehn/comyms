package main

import (
	"log"

	"github.com/todlehn/comyms/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
