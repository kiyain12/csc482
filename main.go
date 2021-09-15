package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	loggly "github.com/jamespearly/loggly"
)

type weatherData struct {
	ID            int     `json:"id"`
	DATE          string  `json:"applicable_date"`
	TEMPERATURE   float64 `json:"the_temp`
	WEATHERSTATUS string  `json:"weather_state_name"`
}

func main() {

	var tag string = "firstapplication"

	client := loggly.New(tag)

	resp, err := http.Get("https://www.metaweather.com/api/location/2459115/2021/9/13/")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var x []weatherData
	err = json.Unmarshal(body, &x)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v", x)

	// _, err = os.Stdout.Write(body)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// jsonFile, err := os.Open(resp.Status)

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println("successfully opened json file")

	// defer jsonFile.Close()

	output := strconv.Itoa(int(len(x)))

	// Valid Send (no error returned)
	err = client.EchoSend("info", "Success! Data size: " + output)
	fmt.Println("err:", err)

}