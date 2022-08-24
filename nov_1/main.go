package main

import (
	"fmt"
	"time"

	natsConnector "nats.pub.net/connector"
	sensorsController "nats.pub.net/sensor"
)

func main() {
	tick := time.Tick(2 * time.Second)

	sensorsController.AddSensor("sensor1")
	sensorsController.AddSensor("sensor2")
	sensorsController.AddSensor("sensor4")

	var sensors []sensorsController.Sensor

	natsConnector.Start("127.0.0.1:4222")

	for range tick {
		newFunction(sensors)
	}
	newFunction(sensors)

	natsConnector.Stop()
	fmt.Println("done")
	defer natsConnector.Stop()
}

func newFunction(sensors []sensorsController.Sensor) {
	sensors = sensorsController.ReadValuesFromSensors()
	natsConnector.SendMessagesFromSensors(sensors)
}
