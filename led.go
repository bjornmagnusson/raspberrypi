package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
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
	mode     = 2
	demoMode = false
	isPushoverEnabled = false
	num 		 = 0
	api 		 = false
	buttonEnabled = false

	// LEDs
	ledRedPin    = 4
	ledYellowPin = 17
	ledGreenPin  = 27
	ledToColor   = map[int]string{0: "red", 1: "yellow", 2: "green"}
	ledMap = map[int]gpio.PinIO{}
	ledMode 		 = 0

	// Buttons
	buttonPin 	= 22
	buttons 		= map[int]gpio.PinIO{}

	// GPIOs
	gpios = map[int]Gpio{}
	demoNum = 26

	// Pushover
	pushoverUser = ""
	pushoverToken = ""
	pushoverApi = "https://api.pushover.net:443/1/messages.json"
)

func getLEDString(color string) string {
	return "Toggle " + strings.ToUpper(color)
}

func buildGpio(id int, name string, value int) Gpio {
	return Gpio{id, name, value}
}

func toggleLED(id int, pin gpio.PinIO, color string) {
	//fmt.Println(getLEDString(color))
	if !demoMode {
		//fmt.Println("Val to write", !pin.Read())
		pin.Out(!pin.Read())
		value := 0
		if pin.Read() == gpio.High {
			value = 1
		}
		gpios[id] = buildGpio(id, pin.Name(), value)
	}
}

func initButtons() {
	buttons[0] = gpioreg.ByName(strconv.Itoa(buttonPin))
	if err := buttons[0].In(gpio.PullUp, gpio.BothEdges); err != nil {
  	log.Fatal(err)
  }
	for buttonIndex := 0; buttonIndex < len(buttons); buttonIndex++ {
		fmt.Printf("%s: %s\n", buttons[buttonIndex], buttons[buttonIndex].Function())
	}
}

func initLEDs() {
	if mode == 2 {
		ledMap[0] = gpioreg.ByName(strconv.Itoa(ledRedPin))
		ledMap[1] = gpioreg.ByName(strconv.Itoa(ledYellowPin))
		ledMap[2] = gpioreg.ByName(strconv.Itoa(ledGreenPin))
		for i := 0; i < len(ledMap); i++ {
			fmt.Println("Resetting ", ledMap[i])
			if err := ledMap[i].Out(gpio.Low); err != nil {
        log.Fatal(err)
      }
		}
	}
}

func initGPIO() error {
	_,err := host.Init()
	return err
}

func doLedToggling(i int, isSleepEnabled bool) {
	if !demoMode {
	  if mode == 2{
			toggleLED(i%3, ledMap[i%3], ledToColor[i%3])
		}
	} else {
		fmt.Println(getLEDString(ledToColor[i%3]))
		gpios[i] = buildGpio(i%demoNum, "GPIO" + strconv.Itoa(i%demoNum), i % 2)
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

func getLedMode(currentLedMode int) int {
	if currentLedMode == 0 {
		return 1
	} else {
		return 0
	}
}

func toggleLedMode() {
	ledMode = getLedMode(ledMode)
}

func sendPushoverMessage(message string) {
	if isPushoverEnabled {
		message := PushoverMessage{pushoverToken, pushoverUser, message}
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

func ledModeHandler(w http.ResponseWriter, r *http.Request) {
	toggleLedMode()
	sendPushoverMessage("LED mode toggled from API")
}

func initWebServer() {
	fmt.Println("Initializing Webserver")
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/gpios", getGpios)
	mux.HandleFunc("/v1/ledMode", ledModeHandler)
	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":8080", handler)
}

func listenForButtonPress(button gpio.PinIO) {
	buttonState := gpio.High
	for {
		fmt.Println("Check button ", button)
		button.WaitForEdge(-1)
		buttonState = button.Read()
		fmt.Printf("-> %s\n", buttonState)
		if buttonState == gpio.High {
			toggleLedMode()
			sendPushoverMessage("LED mode toggled using button")
		}
	}
}

func listenForButtonsPress() {
	fmt.Println("Listening for button presses")
	go listenForButtonPress(buttons[0])
}

func parseParameters() {
	fmt.Println("Parsing parameters")
	numFromCli := flag.Int("num", 0, "number of blinks")
	buttonEnabledFromCli := flag.Bool("button", true, "button mode")
	apiFromCli := flag.Bool("api", true, "API enabled")
	demoFromCli := flag.Bool("demo", false, "Demo mode enabled")
	pushoverFromCli := flag.Bool("pushover", false, "Pushover notifications enabled")

	flag.Parse()

	buttonEnabled = *buttonEnabledFromCli
	api = *apiFromCli
	demoMode = *demoFromCli
	isPushoverEnabled = *pushoverFromCli
	num = *numFromCli
	pushoverUser = os.Getenv("PUSHOVER_USER")
	pushoverToken = os.Getenv("PUSHOVER_TOKEN")
}

func startLedToggling() {
	if num == 0 {
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
		for i := 0; i < num; i++ {
			doLedToggling(i, true)
		}
	}
}

func initPins() {
	if !demoMode {
		var err = initGPIO()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		initLEDs()
		initButtons()
		if buttonEnabled {
			listenForButtonsPress()
		}
	}
}

func main() {
	parseParameters()

	if isPushoverEnabled && (pushoverUser == "" || pushoverToken == "") {
		fmt.Println("Pushover env variables are undefined, disabling pushover notifs")
		isPushoverEnabled = false
	}

	if demoMode {
		fmt.Println("Running in demo mode, no physical hardware interaction")
	}

	fmt.Println("Number of blinks:", num)

	if api {
		go initWebServer()
	}

	initPins()
	startLedToggling()
}
