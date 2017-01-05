package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/all"
	"github.com/stianeikeland/go-rpio"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	mode = 0 // 0=go-rpio,1=embd

	// LEDs
	ledRedPin    = 4
	ledYellowPin = 17
	ledGreenPin  = 27
	ledToColor   = map[int]string{}
	ledMapEmbd = map[int]embd.DigitalPin{}
	ledMap     = map[int]rpio.Pin{}

	// Buttons
	buttonPin = 22
	buttonMap = map[int]embd.DigitalPin{}
)

func getLEDString(color string) string {
	return "Toggle " + strings.ToUpper(color)
}

func getToggledValue(pin embd.DigitalPin) int {
	val, _ := pin.Read()
	if val == embd.High {
		return embd.Low
	} else {
		return embd.High
	}
}

func toggleLEDEmbd(pin embd.DigitalPin, color string) {
	fmt.Println(getLEDString(color))
	toggledValue := getToggledValue(pin)
	fmt.Println("Val to write", toggledValue)
	embd.DigitalWrite(pin.N(), toggledValue)
}

func toggleLED(pin rpio.Pin, color string) {
	fmt.Println(getLEDString(color))
	fmt.Println("Current value:", pin.Read())
	pin.Toggle()
}

func initButtons() {
	buttonMap[0],_ = embd.NewDigitalPin(buttonPin)
	buttonMap[0].SetDirection(embd.In)
	buttonMap[0].ActiveLow(false)
	quit := make(chan interface{})
	err := buttonMap[0].Watch(embd.EdgeFalling, func(btn embd.DigitalPin) {
		fmt.Println("Pressed", btn)
		quit <- btn
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Button %v was pressed.\n", <-quit)
}

func initLEDs() {
	if mode == 1 {
		ledMapEmbd[0],_ = embd.NewDigitalPin(ledRedPin)
		ledMapEmbd[1],_ = embd.NewDigitalPin(ledYellowPin)
		ledMapEmbd[2],_ = embd.NewDigitalPin(ledGreenPin)
		for i := 0; i < len(ledMapEmbd); i++ {
			ledMapEmbd[i].SetDirection(embd.Out)
			embd.DigitalWrite(ledMapEmbd[i].N(), embd.Low)
		}
	} else {
		ledMap[0] = rpio.Pin(ledRedPin)
		ledMap[1] = rpio.Pin(ledYellowPin)
		ledMap[2] = rpio.Pin(ledGreenPin)
		for i := 0; i < len(ledMap); i++ {
			ledMap[i].Output()
			ledMap[i].Low()
		}
	}
	ledToColor[0] = "red"
	ledToColor[1] = "yellow"
	ledToColor[2] = "green"
}

func initGPIO() error {
	if mode == 1 {
		return embd.InitGPIO()
	} else {
		return rpio.Open()
	}
}

func doLedToggling(i int) {
	if mode == 1 {
		toggleLEDEmbd(ledMapEmbd[i%3], ledToColor[i%3])
	} else {
		toggleLED(ledMap[i%3], ledToColor[i%3])
	}
	time.Sleep(time.Second)
}

type Gpio struct {
	ID int
	Name string
	Value int
}

func gpios (w http.ResponseWriter, r *http.Request) {

    m := Gpio{1, "GPIO2", 1}
    b, err := json.Marshal(m)

    if err != nil {
        panic(err)
    }

     w.Write(b)
}

func main() {
	fmt.Println("Parsing parameters")
	num := flag.Int("num", 3, "number of blinks")
	modeFromCli := flag.Int("mode", 0, "mode")
	button := flag.Bool("button", false, "button mode")
	flag.Parse()
	mode = *modeFromCli
	if mode == 1 {
		fmt.Println("Running using Embd.io")
	} else {
		fmt.Println("Running using go-rpio")
	}
	fmt.Println("Number of blinks:", *num)

	http.HandleFunc("/gpios", gpios)
    http.ListenAndServe(":8080", nil)

	var err = initGPIO()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	initLEDs()
	if mode == 1 && *button {
		initButtons()
	}

	if mode == 1 {
		defer embd.CloseGPIO()
	} else {
		defer rpio.Close()
	}

	if !*button {
		if *num == 0 {
			var counter = 0
			for {
				doLedToggling(counter)
				counter++
			}
		} else {
			for i := 0; i < *num; i++ {
				doLedToggling(i)
			}
		}
	}
}
