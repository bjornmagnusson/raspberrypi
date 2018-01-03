package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
	"github.com/rs/cors"
)

var (
	mode     = 2 // 0=go-rpio,2=periph
	demoMode = false
	isPushoverEnabled = false

	// LEDs
	ledRedPin    = 4
	ledYellowPin = 17
	ledGreenPin  = 27
	ledToColor   = map[int]string{}
	ledMapPeriph = map[int]gpio.PinIO{}
	ledMode 		 = 0

	// Buttons
	buttonPin 	= 22
	buttons 		= map[int]gpio.PinIO{}

	// GPIOs
	gpios = map[int]Gpio{}
	demoNum = 26

	// PUSHOVER_USER
	pushoverUser = ""
	pushoverToken = ""
	pushoverApi = "https://api.pushover.net:443/1/messages.json"
)

func getLEDString(color string) string {
	return "Toggle " + strings.ToUpper(color)
}

func setGpio(id int, name string, value int) {
	gpios[id] = Gpio{id, name, value}
}

func toggleLEDPeriph(id int, pin gpio.PinIO, color string) {
	fmt.Println(getLEDString(color))
	if !demoMode {
		fmt.Println("Val to write", !pin.Read())
		pin.Out(!pin.Read())
		value := 0
		if pin.Read() == gpio.High {
			value = 1
		}
		setGpio(id, pin.Name(), value)
	}
}

func initButtons() {
	buttons[0] = gpioreg.ByName(strconv.Itoa(buttonPin))
	fmt.Printf("%s: %s\n", buttons[0], buttons[0].Function())

	buttons[0].In(gpio.PullDown, gpio.BothEdges)
}

func initLEDs() {
	if mode == 2 {
		ledMapPeriph[0] = gpioreg.ByName(strconv.Itoa(ledRedPin))
		ledMapPeriph[1] = gpioreg.ByName(strconv.Itoa(ledYellowPin))
		ledMapPeriph[2] = gpioreg.ByName(strconv.Itoa(ledGreenPin))
		for i := 0; i < len(ledMapPeriph); i++ {
			ledMapPeriph[i].Out(gpio.Low)
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
	_,err := host.Init()
	return err
}

func doLedToggling(i int, isSleepEnabled bool) {
	if !demoMode {
	  if mode == 2{
			toggleLEDPeriph(i%3, ledMapPeriph[i%3], ledToColor[i%3])
		}
	} else {
		fmt.Println(getLEDString(ledToColor[i%3]))
		setGpio(i%demoNum, "GPIO" + strconv.Itoa(i%demoNum), i % 2)
	}

	if isSleepEnabled {
		time.Sleep(time.Second)
	}
}

type Gpio struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func getGpios(w http.ResponseWriter, r *http.Request) {
	gpiosToJson := make([]Gpio, 0, len(gpios))
	fmt.Printf("Transforming %d GPIOs\n", len(gpios))
	for  _, value := range gpios {
  	gpiosToJson = append(gpiosToJson, value)
	}
	fmt.Printf("Transformed into %d gpios\n", len(gpiosToJson))
	json, err := json.Marshal(gpiosToJson)

	if err != nil {
		panic(err)
	}
	w.Write(json)
}

type PushoverMessage struct {
	Token string `json:"token"`
	User string `json:"user"`
	Message string `json:"message"`
}

func ledModeHandler(w http.ResponseWriter, r *http.Request) {
	if ledMode == 0 {
		ledMode = 1
	} else {
		ledMode = 0
	}

	if isPushoverEnabled {
		message := PushoverMessage{pushoverToken, pushoverUser, "LED mode toggled"}
		json, err := json.Marshal(message)
		var jsonStr = []byte(json)

		// Build the request
		client := &http.Client{}
		req, err := http.NewRequest("POST", pushoverApi, bytes.NewBuffer(jsonStr))
		if err != nil {
			fmt.Println("Request ERROR:", err)
			return
		}

	  req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Response ERROR:", err)
			return
		}
		defer resp.Body.Close()
		fmt.Println("Response: ", *resp)
	}
}

func initWebServer() {
	fmt.Println("Initializing Webserver")
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/gpios", getGpios)
	mux.HandleFunc("/v1/ledMode", ledModeHandler)
	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":8080", handler)
}

func listenForButtonPress() {
	fmt.Println("Listening for button presses")
	for {
		for button := 0; button < len(buttons); button++ {
			buttons[button].WaitForEdge(-1)
			fmt.Printf("-> %s\n", buttons[button].Read())
		}
	}
}

func main() {
	fmt.Println("Parsing parameters")
	num := flag.Int("num", 0, "number of blinks")
	buttonEnabled := flag.Bool("button", false, "button mode")
	api := flag.Bool("api", true, "API enabled")
	demo := flag.Bool("demo", false, "Demo mode enabled")
	pushoverFromCli := flag.Bool("pushover", false, "Pushover notifications enabled")

	flag.Parse()

	demoMode = *demo
	isPushoverEnabled = *pushoverFromCli
	pushoverUser = os.Getenv("PUSHOVER_USER")
	pushoverToken = os.Getenv("PUSHOVER_TOKEN")

	if isPushoverEnabled && (pushoverUser == "" || pushoverToken == "") {
		fmt.Println("Pushover env variables are undefined, disabling pushover notifs")
		isPushoverEnabled = false
	}

	if demoMode {
		fmt.Println("Running in demo mode, no physical hw interaction")
	}
	if mode == 2 {
		fmt.Println("Running using periph")
	} else {
		fmt.Println("Running using periph")
		mode = 2
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
		initButtons()
		if *buttonEnabled {
			go listenForButtonPress()
		}
	} else {
		initLEDcolors()
	}

	if *num == 0 {
		var counter = 0
		for {
			if ledMode == 0 {
				doLedToggling(counter, true)
			} else if ledMode == 1 {
				doLedToggling(counter, false)
				doLedToggling(counter + 1, false)
				doLedToggling(counter + 2, true)
			}
			counter++
		}
	} else {
		for i := 0; i < *num; i++ {
			doLedToggling(i, true)
		}
	}
}
