package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio"
)

var pin = rpio.Pin(4)

func main() {
	fmt.Println("Parsing parameters")
	num := flag.Int("num", 0, "number of blinks")
	flag.Parse()

	fmt.Println("Number of blinks: ", *num)

	fmt.Println("Opening rpio access")

	var err = rpio.Open()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer rpio.Close()

	fmt.Println("Pin as output")
	pin.Output()

	for i := 0; i < *num; i++ {
		fmt.Println("Toggle")
		pin.Toggle()
		time.Sleep(time.Second)
	}
}
