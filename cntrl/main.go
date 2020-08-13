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

var adaptEndpoint string = fmt.Sprintf("http://%s:%s/prodcntrl", os.Getenv("ADAPT_IP"), os.Getenv("ADAPT_PORT"))

// standard production rate per second
const rate = 10

func update(queue <-chan struct{}) {
	// update the world every 1 second
	ticker := time.NewTicker(time.Duration(1000) * time.Millisecond)

	discarded := 0
	for {
		select {
		case <-queue:
			discarded++
		case <-ticker.C:
			go func() {
				curr := rate - discarded
				discarded = 0

				id, err := uuid.NewRandom()
				if err != nil {
					log.Print(err)
					return
				}

				type ProdCntrlData struct {
					ProdRate int    `json:"prod_rate"`
					UUID     string `json:"uuid"`
				}

				data, err := json.Marshal(ProdCntrlData{
					ProdRate: curr,
					UUID:     id.String(),
				})

				if err != nil {
					return
				}

				log.Printf("Production in the last second was %d", curr)
				log.Printf("send,prodctrl,%s,%s", id.String(), strconv.FormatInt(time.Now().UnixNano(), 10))

				req, err := http.NewRequest("POST", adaptEndpoint, bytes.NewReader(data))

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
}

func main() {
	discardqueue := make(chan struct{})

	type Request struct {
		UUID string `json:"uuid"`
	}

	http.HandleFunc("/discard", func(w http.ResponseWriter, r *http.Request) {
		timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)

		var data Request
		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			return
		}

		log.Printf("recv,discard,%s,%s", data.UUID, timestamp)

		discardqueue <- struct{}{}
	})

	go update(discardqueue)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("CNTRL_PORT"), nil))

}
