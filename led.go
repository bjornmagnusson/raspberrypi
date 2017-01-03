package main

import (
	"flag"
	"fmt"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/all"
	//"os"
	"strings"
	"time"
	//"github.com/stianeikeland/go-rpio"
)

// var (
// 	ledRed = rpio.Pin(4)
// 	ledYellow = rpio.Pin(17)
// 	ledGreen = rpio.Pin(27)
// )

var (
	ledRed, _    = embd.NewDigitalPin(4)
	ledYellow, _ = embd.NewDigitalPin(17)
	ledGreen, _  = embd.NewDigitalPin(27)
)

func getLEDString(color string) string {
	return "Toggle " + strings.ToUpper(color)
}

func getToggledValue(pin embd.DigitalPin) int {
	val,_ := pin.Read()
	if val == embd.High {
		return embd.Low
	} else {
		return embd.High
	}
}

func toggleLED(pin embd.DigitalPin, color string) {
	fmt.Println(getLEDString(color))
	pin.Write(getToggledValue(pin))
}

// func toggleLED(pin rpio.Pin, color string)  {
// 	fmt.Println(getLEDString(color))
// 	pin.Toggle()
// }

func initLEDs() {
	embd.SetDirection(4, embd.Out)
	embd.SetDirection(17, embd.Out)
	embd.SetDirection(27, embd.Out)
	// ledRed.Output()
	// ledYellow.Output()
	// ledGreen.Output()
}

func main() {
	fmt.Println("Parsing parameters")
	num := flag.Int("num", 0, "number of blinks")
	flag.Parse()

	fmt.Println("Number of blinks:", *num)

	fmt.Println("Opening rpio access")

	//var err = rpio.Open()
	embd.InitGPIO()

	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	//defer rpio.Close()
	defer embd.CloseGPIO()

	fmt.Println("Pin as output")

	for i := 0; i < *num; i++ {
		if i%3 == 0 {
			toggleLED(ledRed, "red")
		} else if i%3 == 1 {
			toggleLED(ledYellow, "yellow")
		} else {
			toggleLED(ledGreen, "green")
		}
		time.Sleep(time.Second)
	}
}
