package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/stianeikeland/go-rpio"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/all"
	"github.com/rs/cors"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

var (
	mode     = 0 // 0=go-rpio,1=embd,2=periph
	demoMode = false

	// LEDs
	ledRedPin    = 4
	ledYellowPin = 17
	ledGreenPin  = 27
	ledToColor   = map[int]string{}
	ledMapEmbd   = map[int]embd.DigitalPin{}
	ledMap       = map[int]rpio.Pin{}
	ledMapPeriph = map[int]gpio.PinIO{}

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
	}
	return embd.High
}

func toggleLEDPeriph(pin gpio.PinIO, color string) {
	fmt.Println(getLEDString(color))
	if !demoMode {
		fmt.Println("Val to write", !pin.Read())
		pin.Out(!pin.Read())
	}
}

func toggleLEDEmbd(pin embd.DigitalPin, color string) {
	fmt.Println(getLEDString(color))
	if !demoMode {
		toggledValue := getToggledValue(pin)
		fmt.Println("Val to write", toggledValue)
		embd.DigitalWrite(pin.N(), toggledValue)
	}
}

func toggleLED(pin rpio.Pin, color string) {
	fmt.Println(getLEDString(color))
	if !demoMode {
		fmt.Println("Current value:", pin.Read())
		pin.Toggle()
	}
}

func initButtons() {
	buttonMap[0], _ = embd.NewDigitalPin(buttonPin)
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
		ledMapEmbd[0], _ = embd.NewDigitalPin(ledRedPin)
		ledMapEmbd[1], _ = embd.NewDigitalPin(ledYellowPin)
		ledMapEmbd[2], _ = embd.NewDigitalPin(ledGreenPin)
		for i := 0; i < len(ledMapEmbd); i++ {
			ledMapEmbd[i].SetDirection(embd.Out)
			embd.DigitalWrite(ledMapEmbd[i].N(), embd.Low)
		}
	} else if mode == 2 {
		ledMapPeriph[0] = gpioreg.ByNumber(ledRedPin)
		ledMapPeriph[1] = gpioreg.ByNumber(ledYellowPin)
		ledMapPeriph[2] = gpioreg.ByNumber(ledGreenPin)
		for i := 0; i < len(ledMapPeriph); i++ {
			ledMapPeriph[i].Out(gpio.Low)
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
	initLEDcolors()
}

func initLEDcolors() {
	ledToColor[0] = "red"
	ledToColor[1] = "yellow"
	ledToColor[2] = "green"
}

func initGPIO() error {
	if mode == 1 {
		return embd.InitGPIO()
	} else if mode == 2 {
		_,err := host.Init()
		return err
	}
	return rpio.Open()
}

func doLedToggling(i int) {
	if !demoMode {
		if mode == 1 {
			toggleLEDEmbd(ledMapEmbd[i%3], ledToColor[i%3])
		} else if mode == 2{
			toggleLEDPeriph(ledMapPeriph[i%3], ledToColor[i%3])
		} else {
			toggleLED(ledMap[i%3], ledToColor[i%3])
		}
	} else {
		fmt.Println(getLEDString(ledToColor[i%3]))
	}
	time.Sleep(time.Second)
}

type Gpio struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func gpios(w http.ResponseWriter, r *http.Request) {
	pins := ledMapEmbd
	numPins := len(pins)
	if demoMode {
		numPins = 26
	}
	fmt.Printf("Fetching %d GPIOs\n", numPins)
	gpios := make([]Gpio, numPins)
	for i := 0; i < numPins; i++ {
		fmt.Println("GPIO", i)
		pin := pins[i]
		if demoMode {
			pinValue := i % 2
			gpio := Gpio{i, "GPIO" + strconv.Itoa(i), pinValue}
			gpios[i] = gpio
		} else {
			pinValue, _ := pin.Read()
			gpio := Gpio{i, "GPIO" + strconv.Itoa(pin.N()), pinValue}
			gpios[i] = gpio
		}
	}
	json, err := json.Marshal(gpios)

	if err != nil {
		panic(err)
	}

	w.Write(json)
}

func initWebServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/gpios", gpios)
	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":8080", handler)
}

func main() {
	fmt.Println("Parsing parameters")
	num := flag.Int("num", 3, "number of blinks")
	modeFromCli := flag.Int("mode", 0, "mode")
	button := flag.Bool("button", false, "button mode")
	api := flag.Bool("api", true, "API enabled")
	demo := flag.Bool("demo", false, "Demo mode enabled")
	flag.Parse()

	mode = *modeFromCli
	demoMode = *demo

	if demoMode {
		fmt.Println("Running in demo mode, no physical hw interaction")
	}
	if mode == 1 {
		fmt.Println("Running using Embd.io")
	} else if mode == 2 {
		fmt.Println("Running using periph")
	} else {
		fmt.Println("Running using go-rpio")
	}
	fmt.Println("Number of blinks:", *num)

	if *api {
		go initWebServer()
	}

	if !demoMode {
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
	} else {
		initLEDcolors()
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
