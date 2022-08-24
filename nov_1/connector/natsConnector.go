package natsConnector

import (
	"fmt"

	nats "github.com/nats-io/nats.go"
	message "nats.pub.net/connector/message"
	sensorsController "nats.pub.net/sensor"
)

// Exported global variable to hold the nats connection pool.

var c *nats.EncodedConn

func Start(addr string) {
	var nc *nats.Conn
	nc, _ = nats.Connect(addr)
	c, _ = nats.NewEncodedConn(nc, nats.JSON_ENCODER)
}

func Send(message message.MessageNats) {

	if err := c.Publish("messages", message); err != nil {
		fmt.Printf("%+v\n", err)
	}

}

func SendMessagesFromSensors(ss []sensorsController.Sensor) {
	var natsMessages []message.SensorMessageNats
	for _, s := range ss {
		natsMessages = append(natsMessages, message.NewSensorMessage(s.Name, s.Timestamp.String(), s.Value))
	}

	if err := c.Publish("messages", message.NewMessage(natsMessages)); err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func Stop() {
	c.Close()
}
