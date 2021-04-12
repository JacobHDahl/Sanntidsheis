package main

import (
	"flag"
	"fmt"

	"./FSM"
	"./config"
	"./driver/elevio"
	"./network/peers"
)

const numFloors = 4

func main() {
	elevPort_p := flag.String("elev_port", "15657", "The port of the elevator to connect to (for sim purposes)")

	flag.Parse()

	elevPort := *elevPort_p
	hostString := "localhost:" + elevPort

	fmt.Println("Elevport ", hostString)

	println("Connecting to server")
	elevio.Init(hostString, numFloors)

	driverChannels := config.DriverChannels{
		DrvButtons:     make(chan elevio.ButtonEvent),
		DrvFloors:      make(chan int),
		DrvStop:        make(chan bool),
		DoorsOpen:      make(chan int),
		CompletedOrder: make(chan elevio.ButtonEvent, 100),
		DrvObstr:       make(chan bool),
	}

	dummyString := "fuck u bitch"
	transmitEnable := make(chan bool)
	go peers.Transmitter(22349, dummyString, transmitEnable)

	go elevio.PollObstructionSwitch(driverChannels.DrvObstr)
	go elevio.PollButtons(driverChannels.DrvButtons)
	go elevio.PollFloorSensor(driverChannels.DrvFloors)
	go elevio.PollStopButton(driverChannels.DrvStop)
	go FSM.Fsm(driverChannels.DoorsOpen)

	FSM.InternalControl(driverChannels)

}
