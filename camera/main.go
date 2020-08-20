package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

var cfdEndpoint string = fmt.Sprintf("http://%s:%s/image", os.Getenv("CFD_IP"), os.Getenv("CFD_PORT"))

// production rate defined in cntrl
const rate = 10

const width int = 100
const height int = 100

func generateImage() string {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if rand.Intn(2) == 1 {
				img.Set(x, y, color.White)
			} else {
				img.Set(x, y, color.Black)
			}
		}
	}

	var r bytes.Buffer

	err := png.Encode(&r, img)

	if err != nil {
		return ""
	}

	reader := bufio.NewReader(&r)
	content, _ := ioutil.ReadAll(reader)

	// Encode as base64.
	encoded := base64.StdEncoding.EncodeToString(content)

	return encoded
}

func generateImages() {
	ticker := time.NewTicker(time.Duration(1.0 / rate * 1000.0) * time.Millisecond)

	for range ticker.C {

		img := generateImage()

		id, err := uuid.NewRandom()
		if err != nil {
			log.Print(err)
			continue
		}

		type Request struct {
			Img  string `json:"img"`
			UUID string `json:"uuid"`
		}

		data, err := json.Marshal(Request{
			Img:  img,
			UUID: id.String(),
		})

		if err != nil {
			return
		}

		log.Printf("send,camera,%s,%s", id.String(), strconv.FormatInt(time.Now().UnixNano(), 10))

		req, err := http.NewRequest("POST", cfdEndpoint, bytes.NewReader(data))

		if err != nil {
			continue
		}

		go func() {
			_, err = (&http.Client{}).Do(req)

			if err != nil {
				log.Print(err)

			}
		}()

	}
}

func main() {
	rand.Seed(100)

	http.HandleFunc("/state/notifications", func(w http.ResponseWriter, r *http.Request) {
		stateNameA, ok := r.URL.Query()["state_name"]

		if !ok || len(stateNameA[0]) < 1 {
			log.Println("Url Param 'state_name' is missing")
			return
		}

		stateName := stateNameA[0]

		log.Printf("Received state notification request, beginning " + string(stateName))
		rand.Seed(100)
		w.Header().Set("Server", "camera")
		w.WriteHeader(200)
	})

	go generateImages()
	
	log.Printf("started")

	log.Fatal(http.ListenAndServe(":"+os.Getenv("CAMERA_PORT"), nil))
}
