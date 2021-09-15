package main

import (
	"fmt"
	loggly "github.com/jamespearly/loggly"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

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

	// Valid Send (no error returned)
	err = client.Send("info", body)
	fmt.Println("err:", err)

}
