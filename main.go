package main

import (
	"fmt"
	"time"

	"../Sanntidsheis/FSM"
	"../Sanntidsheis/config"
	"../Sanntidsheis/driver/elevio"
)

func main() {

	driverChannels := config.DriverChannels{
		DrvButtons:     make(chan elevio.ButtonEvent),
		DrvFloors:      make(chan int),
		DrvStop:        make(chan bool),
		DoorsOpen:      make(chan int),
		CompletedOrder: make(chan elevio.ButtonEvent, 100),
	}

	fmt.Println("test", time.Second)
	FSM.InternalControl(driverChannels)

}
