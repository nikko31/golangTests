package sensor

import (
	"math/rand"
	"strconv"
	"time"
)

type Sensor struct {
	Name      string
	Timestamp time.Time
	Value     string
}

var sensorsList []Sensor

func AddSensor(name string) {
	sensor := Sensor{Name: name}
	sensorsList = append(sensorsList, sensor)
}

func updateValuesOfSensor(name string) {
	for i := 0; i < len(sensorsList); i++ {
		sensor := &sensorsList[i]
		if sensor.Name == name {
			sensor.Timestamp = time.Now()
			sensor.Value = strconv.Itoa(rand.Intn(100))
		}
	}
}

func updateValues() {
	for _, sensor := range sensorsList {
		updateValuesOfSensor(sensor.Name)
	}
}

func ReadValuesFromSensors() []Sensor {
	updateValues()
	return sensorsList
}
