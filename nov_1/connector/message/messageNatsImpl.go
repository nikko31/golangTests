package message

type SensorMessageNats struct {
	Name      string `json:"Name"`
	Timestamp string `json:"Timestamp"`
	Value     string `json:"Value"`
}

type MessageNats struct {
	Sensor []SensorMessageNats `json:"sensor"`
}

func NewSensorMessage(name string, timestamp string, value string) SensorMessageNats {
	var sensor = SensorMessageNats{name, timestamp, value}
	return sensor
}

func NewMessage(sensorToAdd []SensorMessageNats) MessageNats {
	var message = MessageNats{sensorToAdd}
	return message
}
