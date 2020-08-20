package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

var adaptEndpoint string = fmt.Sprintf("http://%s:%s/sensor", os.Getenv("ADAPT_IP"), os.Getenv("ADAPT_PORT"))

const interval int = 100
const upper int = 150
var lower int = 50

func generateTemps() {
	// update the world every 100 milliseconds
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	for range ticker.C {
		go func() {
			temp := rand.Intn(upper-lower) + lower
			// send random temperature

			id, err := uuid.NewRandom()
			if err != nil {
				log.Print(err)
				return
			}

			type Reading struct {
				Temp int    `json:"temp"`
				UUID string `json:"uuid"`
			}

			data, err := json.Marshal(Reading{
				Temp: temp,
				UUID: id.String(),
			})

			if err != nil {
				return
			}

			log.Printf("send,adapt,%s,%s", id.String(), strconv.FormatInt(time.Now().UnixNano(), 10))

			go func() {
				req, err := http.NewRequest("POST", adaptEndpoint, bytes.NewReader(data))

				if err != nil {
					return
				}
				_, err = (&http.Client{}).Do(req)

				if err != nil {
					log.Print(err)
				}
			}()
		}()
	}
}

func main() {
	rand.Seed(100)

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {

		minTempA, ok := r.URL.Query()["min_temp"]

		if !ok || len(minTempA[0]) < 1 {
			log.Println("Url Param 'min_temp' is missing")
			return
		}

		minTemp, err := strconv.Atoi(minTempA[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		log.Printf("Received new minimum temperature, was " + strconv.Itoa(lower) + ", is now " + strconv.Itoa(minTemp))

		lower = minTemp

		w.Header().Set("Server", "temperature-sensor")
    	w.WriteHeader(200)
	})

	http.HandleFunc("/state/notifications", func(w http.ResponseWriter, r *http.Request) {
		stateNameA, ok := r.URL.Query()["state_name"]

		if !ok || len(stateNameA[0]) < 1 {
			log.Println("Url Param 'state_name' is missing")
			return
		}

		stateName := stateNameA[0]

		log.Printf("Received state notification request, beginning " + string(stateName))
		rand.Seed(100)
		w.Header().Set("Server", "temperature-sensor")
		w.WriteHeader(200)
	})

	go generateTemps()
	
	log.Printf("started")

	log.Fatal(http.ListenAndServe(":"+os.Getenv("SENSOR_PORT"), nil))
}
