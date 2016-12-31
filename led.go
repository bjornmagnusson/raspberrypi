package main

import "fmt"
import "github.com/stianeikeland/go-rpio"
import "os"
import "time"

var pin = rpio.Pin(4)

func main() {
    fmt.Println("Opening rpio")

    var err = rpio.Open()

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    defer rpio.Close()

    fmt.Println("Pin as output")
    pin.Output()

    for i := 0; i < 20; i++ {
        fmt.Println("Toggle")
        pin.Toggle()
        time.Sleep(time.Second)
    }
}
