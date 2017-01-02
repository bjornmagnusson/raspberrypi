package main

import (
	"flag"
	"fmt"
	"os"
	"time"
	"strings"
	//"github.com/kidoman/embd"
	"github.com/stianeikeland/go-rpio"
)

var (
	ledRed = rpio.Pin(4)
	ledYellow = rpio.Pin(17)
	ledGreen = rpio.Pin(27)
)

func getLEDString(color string) string {
	return "Toggle" + strings.ToUpper(color)
}

func toggleLED(pin rpio.Pin, color string)  {
	fmt.Println(getLEDString(color))
	pin.Toggle()
}

func initLEDs() {
	ledRed.Output()
	ledYellow.Output()
	ledGreen.Output()
}

func main() {
	fmt.Println("Parsing parameters")
	num := flag.Int("num", 0, "number of blinks")
	flag.Parse()

	fmt.Println("Number of blinks:", *num)

	fmt.Println("Opening rpio access")

	var err = rpio.Open()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer rpio.Close()

	fmt.Println("Pin as output")


	for i := 0; i < *num; i++ {
		if i % 3 == 0 {
			toggleLED(ledRed, "red")
		} else if i % 3 == 1 {
			toggleLED(ledYellow, "yellow")
		} else {
			toggleLED(ledGreen, "green")
		}
		time.Sleep(time.Second)
	}
}
