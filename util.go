package main

import (
	"log"
	"os"
)

func setupLogging() {
	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile("go-tdd.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.Println("Logging initiated.")
}

func chkErr(err error, msg string) {
	if err != nil {
		log.Fatalf("Error: %v, %s", err, msg)
	}
}
