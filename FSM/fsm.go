package main

import (
	"fmt"

	"../driver/elevio"
)

/*
func goToFloor(floorRequest int, elevatorState <-chan int) {
    currentFloor := <- elevatorState

    fmt.Println(currentFloor)

}
*/
func main() {

	//numFloors := 4

	//queue := make([]elevio.ButtonEvent, 0)
	const numFloors = 4
	const numButtons = 4

	var orderMatrix [numFloors][numButtons]int

	elevio.Init("localhost:15657", numFloors)

	var d elevio.MotorDirection = elevio.MD_Up
	//elevio.SetMotorDirection(d)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	type ElevatorState struct {
		IDLE       bool
		RUNNING    bool
		OBSTRUCTED bool
	}

	state := ElevatorState{true, false, false}

	for {
		select {
		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)
			elevio.SetButtonLamp(a.Button, a.Floor, true)

			orderMatrix[int(a.Button)][a.Floor] = 1

		case a := <-drv_floors:
			fmt.Printf("Swag %+v\n", a)
			if a == numFloors-1 {
				d = elevio.MD_Down
			} else if a == 0 {
				d = elevio.MD_Up
			}
			elevio.SetMotorDirection(d)

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				elevio.SetMotorDirection(d)
			}

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < numFloors; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}
		}

	}
}

func calculateMotorDir(floor, order int) bool {
	if order > floor {
		return true
	} else if order < floor {
		return false
	}
	return false
}

func fsm() {

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	//Floor: n Button : n
	//

	type ElevatorState struct {
		IDLE       bool
		RUNNING    bool
		OBSTRUCTED bool
	}

	for {
		currentFloor := <-drv_floors
		currentButton := <-drv_buttons
		currentObstruction := <-drv_obstr
		currentStop := <-drv_stop
		currentOrder := currentButton.Floor

		state := ElevatorState{true, false, false}

		if currentObstruction || currentStop {
			state.OBSTRUCTED = true
		}

		switch {
		case state.IDLE:
			dir := calculateMotorDir(currentFloor, currentOrder)

			if dir {
				elevio.SetMotorDirection(elevio.MD_Up)
			} else {
				elevio.SetMotorDirection(elevio.MD_Down)
			}

			state.RUNNING = true
			state.IDLE = false

		case state.RUNNING:
			if currentOrder == currentFloor {
				elevio.SetMotorDirection(elevio.MD_Stop)
				state.RUNNING = false
				state.IDLE = true
			}

		case state.OBSTRUCTED:
			elevio.SetMotorDirection(elevio.MD_Stop)

		}

	}
}

/*
Dere må dekalrer een matrise [floor][buttons]boool
floor (må finne ut hvileke etajse heisen er i)
    hvis ikk i etasje send til nærnmeste etajse
state_Elevator = idle



deklarer button timer
dekalrer motor kræsj timer

select case
    case tar inn channel order
        lagre den i matrisen

        state machine
            case heise er i idle
                motor = kalkuler_motor__retning(floor, ordre)
                state_Elevator = running
                starte motor_timer
            case heis er i running

            case


    casse channel finis door timer
        slår vi av door lyse
        sjekker og eventuelt kalkulerer nå moor retning
            hvis det er andre order vi har så starter motor og starte mtoor timer
            og setter state til running

        hvis det ikke er sant
            så setter state til idle

    case channel motortimer går ut
        så må vi slå alarm

    case channel floor
        if true er heisen i etajsen den vurder stoppe()
            starte door timer
            slå på door lys
            stoppe motor timer

    case hei er i door










*/