package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	_ "github.com/lib/pq"
	nats "github.com/nats-io/nats.go"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "username"
	password = "password"
	dbname   = "postgres"
)

type SensorMessage struct {
	Name      string `json:"Name"`
	Timestamp string `json:"Timestamp"`
	Value     string `json:"Value"`
}

type Message struct {
	Sensor []SensorMessage `json:"sensor"`
}

var nc *nats.Conn
var c *nats.EncodedConn

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")

	go subStart(db)
	fmt.Scanln()
	defer db.Close()
	defer c.Close()
	nc.Drain()
	// Close connection
	nc.Close()

}

func subStart(db *sql.DB) {

	nc, _ = nats.Connect("127.0.0.1:4222")
	c, _ = nats.NewEncodedConn(nc, nats.JSON_ENCODER)

	c.QueueSubscribe("messages", "queue_titi", func(msg *nats.Msg) {
		var message Message
		err := json.Unmarshal([]byte(msg.Data), &message)
		if err != nil {
			fmt.Printf("%+v\n", err)
		}
		fmt.Printf("%+v\n", "message from pub")
		fmt.Printf("%+v\n", message)

		var sensorIds [3]int
		var avg float32 = 0
		for i, sensor := range message.Sensor {
			insertRowData := `
INSERT INTO SENSORS_DATA (NAME, TIMESTAMP, VALUE)
VALUES ($1, $2, $3)
RETURNING id`
			id := 0
			fmt.Println("New record  is:" + sensor.Name + sensor.Timestamp + sensor.Value)
			err = db.QueryRow(insertRowData, sensor.Name, sensor.Timestamp, sensor.Value).Scan(&id)
			if err != nil {
				panic(err)
			}
			fmt.Println("New record ID is:", id)
			sensorIds[i] = id
			value, _ := strconv.Atoi(sensor.Value)
			avg = avg + float32(value)
		}

		insertAvg := `
INSERT INTO SENSORS_AVG (ID_SENSOR1, ID_SENSOR2, ID_SENSOR3, AVERAGE)
VALUES ($1, $2, $3, $4)
RETURNING id`
		id := 0
		err = db.QueryRow(insertAvg, sensorIds[0], sensorIds[1], sensorIds[2], avg/float32(len(sensorIds))).Scan(&id)
		if err != nil {
			panic(err)
		}
		fmt.Println("New record ID is:", id)

		c.Flush()
	})

}
