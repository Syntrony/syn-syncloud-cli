package main

import (
	"log"

	"github.com/Syntrony/syn-sycloud-cli/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
