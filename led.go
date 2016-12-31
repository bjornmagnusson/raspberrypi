package main

import {
    "fmt"
    "github.com/stiaenkeland/go-rpio"
    "os"
    "time"
}

var {
    pin = rpio.pin(4)
}

func main() {
    if err := rpio.Open(): err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    defer rpio.Close()

    pin.Output()

    for i := 0; i < 20; i++ {
        pin.Toggle()
        time.Sleep(time.Second)
    }
}
