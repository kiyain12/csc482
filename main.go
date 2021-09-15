package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	loggly "github.com/jamespearly/loggly"
	"fmt"
)

func main() {
	resp, err := http.Get("https://www.metaweather.com/api/location/2459115/2021/9/13/")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stdout.Write(body)

	if err != nil {
		log.Fatal(err)
	}

	// jsonFile, err := os.Open(resp.Status)

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println("successfully opened json file")

	// defer jsonFile.Close()

	var tag string = "My-Go-Demo"
	

	// client := loggly.New(tag)

	// log.Print("This is our first log message in Go.")

	client := loggly.New(tag)

		// Valid Send (no error returned)
		err = client.Send("error", "Good morning! No echo.")
		fmt.Println("err:", err)

		
		// Invalid EchoSend -- message level error
		err = client.EchoSend("blah", "blah")
		fmt.Println("err:", err)


}
