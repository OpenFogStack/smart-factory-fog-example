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

	"github.com/google/uuid"
)

var predictEndpoint string = fmt.Sprintf("http://%s:%s/input", os.Getenv("PREDICT_IP"), os.Getenv("PREDICT_PORT"))

var aggregateEndpoint string = fmt.Sprintf("http://%s:%s/data", os.Getenv("AGGREGATE_IP"), os.Getenv("AGGREGATE_PORT"))

// update interval in milliseconds
const interval int = 100

var rate int = 0
var backlog int = 0

func main() {

	type PackCtrlData struct {
		Rate    int    `json:"rate"`
		Backlog int    `json:"backlog"`
		UUID    string `json:"uuid"`
	}

	http.HandleFunc("/rate", func(w http.ResponseWriter, r *http.Request) {
		timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)

		var data PackCtrlData
		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			return
		}

		log.Printf("recv,rate,%s,%s", data.UUID, timestamp)

		rate = data.Rate
		backlog = data.Backlog
	})

	go func() { log.Fatal(http.ListenAndServe(":"+os.Getenv("PKGCNTRL_PORT"), nil)) }()

	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)

	for range ticker.C {
		// both requests need to be performed in parallel, hence we use goroutines for both
		go func() {
			id, err := uuid.NewRandom()
			if err != nil {
				log.Print(err)
				return
			}

			// send rate and backlog
			data, err := json.Marshal(PackCtrlData{
				Rate:    rate,
				Backlog: backlog,
				UUID:    id.String(),
			})

			if err != nil {
				log.Print(err)
				return
			}

			log.Printf("send,packctrl,%s,%s", id.String(), strconv.FormatInt(time.Now().UnixNano(), 10))

			req, err := http.NewRequest("POST", predictEndpoint, bytes.NewReader(data))

			if err != nil {
				return
			}
			_, err = (&http.Client{}).Do(req)

			if err != nil {
				log.Print(err)
			}

			// both requests need to be performed in parallel, hence we use goroutines for both
		}()
		go func() {

			id, err := uuid.NewRandom()
			if err != nil {
				log.Print(err)
				return
			}

			// send rate and backlog
			data, err := json.Marshal(PackCtrlData{
				Rate:    rate,
				Backlog: backlog,
				UUID:    id.String(),
			})

			if err != nil {
				log.Print(err)
				return
			}

			log.Printf("send,packcntrl,%s,%s", id.String(), strconv.FormatInt(time.Now().UnixNano(), 10))

			req, err := http.NewRequest("POST", aggregateEndpoint, bytes.NewReader(data))

			if err != nil {
				return
			}
			_, err = (&http.Client{}).Do(req)

			if err != nil {
				log.Print(err)
			}

		}()
	}
}
