#!/bin/bash
echo "Velkommen til heis!"
echo "Alt er hardkoda nå:)"
cd Simulator-v2
gnome-terminal -x ./SimElevatorServer --port=15555
sleep 1
gnome-terminal -x ./SimElevatorServer --port=15444
cd ..
gnome-terminal -x go run main.go -elev_port=15555 -transmit_port=1111 -receive_port=2222  -elev_id="1"
sleep 1
go run main.go -elev_port=15444 -transmit_port=2222 -receive_port=1111  -elev_id="2"



