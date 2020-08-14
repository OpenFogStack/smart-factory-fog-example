package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const upperthreshold int = 130

// maximum number of items to be packaged per second
const maxpackagingspeed = 15

var packcntrlEndpoint string = fmt.Sprintf("http://%s:%s/rate", os.Getenv("PKGCNTRL_IP"), os.Getenv("PKGCNTRL_PORT"))

type ProdCntrlData struct {
	ProdRate int    `json:"prod_rate"`
	UUID     string `json:"uuid"`
}

type SensorData struct {
	Temp int    `json:"temp"`
	Time string `json:"time"`
	UUID string `json:"uuid"`
}

func update(packagingspeed int, backlog int, id string) {

	type PackCtrlData struct {
		Rate    int    `json:"rate"`
		Backlog int    `json:"backlog"`
		UUID    string `json:"uuid"`
	}

	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)

	// send data
	data, err := json.Marshal(PackCtrlData{
		Rate:    packagingspeed,
		Backlog: backlog,
		UUID:    id,
	})

	if err != nil {
		return
	}

	log.Printf("send,adapt,%s,%s", id, timestamp)

	req, err := http.NewRequest("POST", packcntrlEndpoint, bytes.NewReader(data))

	if err != nil {
		return
	}

	_, err = (&http.Client{}).Do(req)

	if err != nil {
		log.Print(err)
	}
}

func packagingRate(prodrate <-chan ProdCntrlData, temp <-chan SensorData) {
	// backlog of parts that need to be packaged
	var backlog int

	// packaging speed in parts packaged per second
	var packagingspeed int

	// is the temperature blocking packaging at the moment
	var temperatureblock bool

	// production speed of the production machine
	var productionspeed int

	for {
		select {
		case pr := <-prodrate:
			productionspeed = pr.ProdRate

			if temperatureblock {
				packagingspeed = 0
				backlog = backlog + productionspeed
			} else if backlog > 0 {
				packagingspeed = maxpackagingspeed
			} else {
				packagingspeed = productionspeed
			}

			backlog = backlog - packagingspeed
			if backlog < 0 {
				backlog = 0
			}

			go update(packagingspeed, backlog, pr.UUID)

		case t := <-temp:
			temperatureblock = t.Temp > upperthreshold
		}
	}
}

func main() {
	// HTTP service, two endpoints:
	// 1: prod input, change packaging rate
	// 2: sensor input, stop production if threshold is exceeded
	// when temperature returns to normal, process backlog

	proddata := make(chan ProdCntrlData)

	sensordata := make(chan SensorData)

	http.HandleFunc("/prodcntrl", func(w http.ResponseWriter, r *http.Request) {
		timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)

		var data ProdCntrlData
		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			return
		}

		log.Printf("recv,prodcntrl,%s,%s", data.UUID, timestamp)

		proddata <- data
	})

	http.HandleFunc("/sensor", func(w http.ResponseWriter, r *http.Request) {
		timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)

		var data SensorData
		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			return
		}

		log.Printf("recv,sensor,%s,%s", data.UUID, timestamp)

		sensordata <- data
	})

	go packagingRate(proddata, sensordata)

	log.Printf("started")

	log.Fatal(http.ListenAndServe(":"+os.Getenv("ADAPT_PORT"), nil))
}
