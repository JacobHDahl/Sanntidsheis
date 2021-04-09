package FSM

import (
	"fmt"
	"time"

	"../driver/elevio"
)

/*
func goToFloor(floorRequest int, elevatorState <-chan int) {
    currentFloor := <- elevatorState

    fmt.Println(currentFloor)

}
*/

const numFloors = 4
const numButtons = 3

var elevator = Elev{
	State: IDLE,
	Dir:   STILL,
	Floor: 0,
	Queue: [numFloors][numButtons]bool{},
}

var dir elevio.MotorDirection

func FsmInit() {

	elevator.State = IDLE
	// Needs to start in a well-defined state
	for elevator.Floor = elevio.GetFloor(); elevator.Floor < 0; elevator.Floor = elevio.GetFloor() {
		elevio.SetMotorDirection(elevio.MD_Up)
		time.Sleep(1 * time.Second)
	}
	elevio.SetMotorDirection(elevio.MD_Stop)
	fmt.Println("FSM initialized!")
}

func FsmUpdateFloor(newFloor int) {
	elevator.Floor = newFloor
}

func removeButtonLamps(elevator Elev) {
	elevio.SetButtonLamp(elevio.BT_Cab, elevator.Floor, false)
	elevio.SetButtonLamp(elevio.BT_HallDown, elevator.Floor, false)
	elevio.SetButtonLamp(elevio.BT_HallUp, elevator.Floor, false)
}

func fsm(doorsOpen chan<- int) {
	for {
		switch elevator.State {
		case IDLE:
			if ordersAbove(elevator) {
				//println("order above,going up, current Floor: ", Floor)
				dir = elevio.MD_Up
				elevator.Dir = motorDirToElevDir(dir)
				elevio.SetMotorDirection(dir)
				elevator.State = RUNNING
			}
			if ordersBelow(elevator) {
				//println("order below, going down, current Floor: ", Floor)
				dir = elevio.MD_Down
				elevator.Dir = motorDirToElevDir(dir)
				elevio.SetMotorDirection(dir)
				elevator.State = RUNNING
			}
			if ordersInFloor(elevator) {
				//println("order below, going down, current Floor: ", Floor)
				elevator.State = DOOR_OPEN
			}
		case RUNNING:
			if ordersInFloor(elevator) { // this is the problem : the floor is being kept constant at e.g. 2 while its moving
				dir = elevio.MD_Stop
				elevator.Dir = motorDirToElevDir(dir)
				elevio.SetMotorDirection(dir)
				elevator.State = DOOR_OPEN
			}
		case DOOR_OPEN:
			fmt.Println("printing queue")
			printQueue(elevator)
			elevio.SetDoorOpenLamp(true)
			dir = elevio.MD_Stop
			elevio.SetMotorDirection(dir)
			println("DOOR OPEN")
			DeleteOrder(&elevator)
			elevator.State = IDLE
			doorsOpen <- elevator.Floor
			timer1 := time.NewTimer(2 * time.Second)
			<-timer1.C
			elevio.SetDoorOpenLamp(false)
			removeButtonLamps(elevator)
			println("DOOR CLOSE")

		}
	}

}

// InternalControl .. Responsible for internal control of a single elevator
func InternalControl() {
	println("Connecting to server")
	elevio.Init("localhost:15657", numFloors)

	FsmInit()
	drvButtons := make(chan elevio.ButtonEvent)
	drvFloors := make(chan int)
	drvStop := make(chan bool)
	doorsOpen := make(chan int)
	CompletedOrder := make(chan elevio.ButtonEvent, 100)

	//drv_obstr := make(chan bool)
	//go elevio.PollObstructionSwitch(drv_obstr)

	go elevio.PollButtons(drvButtons)
	go elevio.PollFloorSensor(drvFloors)
	go elevio.PollStopButton(drvStop)
	go fsm(doorsOpen)
	for {
		select {
		case floor := <-drvFloors: //Sensor senses a new floor
			//println("updating floor:", floor)
			FsmUpdateFloor(floor)
		case drvOrder := <-drvButtons: // a new button is pressed on this elevator
			//ch.DelegateOrder <- drvOrder //Delegate this order
			fmt.Println("New order")
			fmt.Println(drvOrder)
			elevator.Queue[drvOrder.Floor][int(drvOrder.Button)] = true
			elevio.SetButtonLamp(drvOrder.Button, drvOrder.Floor, true)
		/*case ExtOrder := <-ch.TakeExternalOrder:
		AddOrder(ExtOrder)*/
		case floor := <-doorsOpen:
			order_OutsideUp_Completed := elevio.ButtonEvent{
				Floor:  floor,
				Button: elevio.BT_HallUp,
			}
			order_OutsideDown_Completed := elevio.ButtonEvent{
				Floor:  floor,
				Button: elevio.BT_HallDown,
			}
			order_Inside_Completed := elevio.ButtonEvent{
				Floor:  floor,
				Button: elevio.BT_Cab,
			}
			CompletedOrder <- order_OutsideUp_Completed
			CompletedOrder <- order_OutsideDown_Completed
			CompletedOrder <- order_Inside_Completed

		case <-drvStop:
			elevio.SetMotorDirection(elevio.MD_Stop)
			time.Sleep(3 * time.Second)

		}

	}
}
