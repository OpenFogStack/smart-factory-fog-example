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

const interval int = 10
const upper int = 150
const lower int = 50

func main() {
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
