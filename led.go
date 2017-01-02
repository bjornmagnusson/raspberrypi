package main

import (
	"flag"
	"fmt"
	"os"
	"time"
	//"github.com/kidoman/embd"
	"github.com/stianeikeland/go-rpio"
)

var (
	ledRed = rpio.Pin(4)
	ledYellow = rpio.Pin(17)
	ledGreen = rpio.Pin(27)
)

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
	ledRed.Output()
	ledYellow.Output()
	ledGreen.Output()

	for i := 0; i < *num; i++ {
		if i % 3 == 0 {
			fmt.Println("Toggle RED")
			ledRed.Toggle()
		} else if i % 3 == 1 {
			fmt.Println("Toggle YELLOW")
			ledYellow.Toggle()
		} else {
			fmt.Println("Toggle GREEN")
			ledGreen.Toggle()
		}
		time.Sleep(time.Second)
	}
}
