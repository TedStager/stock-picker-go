package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type chart struct {
	sym Symbol
	date []string
	open, close, high, low []float32
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		log.Fatal(err)
	}
}

func main() {
	dat_chan := make(chan chart)
	go scraper(dat_chan)
	// time.Sleep(2 * time.Second)
	// go scraper(dat_chan)

	hits := picker(dat_chan)
	fmt.Println(hits)

	// write to file
	str_out := strings.Join(hits, "\n")
	err := os.WriteFile("out.dat", []byte(str_out), 777)
	checkError(err)

	fmt.Println("done")
}