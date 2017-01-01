package main

import "fmt"
import "github.com/stianeikeland/go-rpio"
import "os"
import "time"

var pin = rpio.Pin(4)

func main() {
	fmt.Println("Opening rpio access")

	var err = rpio.Open()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer rpio.Close()

	fmt.Println("Pin as output")
	pin.Output()

	fmt.Println("Number of blinks?")
	var num int
	_, err := fmt.Scanf("%d", &num)
	if err != nil {
		fmt.Println("not a number")
	} else {
		fmt.Print(num)
		fmt.Println(" is a number")
	}

	for i := 0; i < num; i++ {
		fmt.Println("Toggle")
		pin.Toggle()
		time.Sleep(time.Second)
	}
}
