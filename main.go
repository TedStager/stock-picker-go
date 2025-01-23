package main

import (
	"fmt"
	"log"
)

func checkError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		log.Fatal(err)
	}
}

func main() {
	scraper()
}