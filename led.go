package main

import "fmt"
import "github.com/stianeikeland/go-rpio"
import "os"
import "time"

var pin = rpio.Pin(4)

func main() {
    rpio.Open()

    defer rpio.Close()

    pin.Output()

    for i := 0; i < 20; i++ {
        pin.Toggle()
        time.Sleep(time.Second)
    }
}
