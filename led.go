package main

import (
	"flag"
	"fmt"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/all"
	"os"
	"strings"
	"time"
	"github.com/stianeikeland/go-rpio"
)

var (
	mode = 0 // 0=go-rpio,1=embd

	ledRedPin = 4
	ledYellowPin = 17
	ledGreenPin = 27
	ledToColor = map[int]string{}

	ledRed = rpio.Pin(4)
	ledYellow = rpio.Pin(17)
	ledGreen = rpio.Pin(27)

	ledRedEmbd, _    = embd.NewDigitalPin(7)
	ledYellowEmbd, _ = embd.NewDigitalPin(0)
	ledGreenEmbd, _  = embd.NewDigitalPin(2)

	ledMapEmbd = map[int]embd.DigitalPin{}
	ledMapRpio = map[int]rpio.Pin{}
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

// func toggleLED(pin embd.DigitalPin, color string) {
// 	fmt.Println(getLEDString(color))
// 	toggledValue := getToggledValue(pin)
// 	fmt.Println("Val to write", toggledValue)
// 	embd.DigitalWrite(pin.N(), embd.High)
// }

func toggleLED(pin rpio.Pin, color string)  {
	fmt.Println(getLEDString(color))
	fmt.Println("Current value:",pin.Read())
	pin.Toggle()
}

func initLEDs() {
	if mode == 1 {
		ledRedEmbd.SetDirection(embd.Out)
		ledYellowEmbd.SetDirection(embd.Out)
		ledGreenEmbd.SetDirection(embd.Out)
	} else {
		ledRed.Output()
		ledYellow.Output()
		ledGreen.Output()
		ledToColor[0] = "red"
		ledToColor[1] = "yellow"
		ledToColor[2] = "green"
	}
}

func main() {
	fmt.Println("Parsing parameters")
	num := flag.Int("num", 0, "number of blinks")
	flag.Parse()

	fmt.Println("Number of blinks:", *num)

	fmt.Println("Opening rpio access")

	var err = rpio.Open()
	// embd.InitGPIO()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	initLEDs()

	if mode == 1 {
		defer embd.CloseGPIO()
	} else {
		defer rpio.Close()
	}

	fmt.Println("Pin as output")

	for i := 0; i < *num; i++ {
		if i%3 == 0 {
			toggleLED(ledRed, ledToColor[0])
		} else if i%3 == 1 {
			toggleLED(ledYellow, ledToColor[1])
		} else {
			toggleLED(ledGreen, ledToColor[2])
		}
		time.Sleep(time.Second)
	}
}
