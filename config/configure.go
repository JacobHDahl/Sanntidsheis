package config

import (
	"../driver/elevio"
)

type DriverChannels struct {
	DrvButtons     chan elevio.ButtonEvent
	DrvFloors      chan int
	DrvStop        chan bool
	DoorsOpen      chan int
	CompletedOrder chan elevio.ButtonEvent
}
