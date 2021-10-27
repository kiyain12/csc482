package main

// https://www.metaweather.com/api/
// https://kiyain.loggly.com/search#terms=&from=2021-09-13T14:54:10.521Z&until=2021-09-15T14:54:10.521Z&source_group=
// https://tutorialedge.net/golang/parsing-json-with-golang/
// https://stackoverflow.com/questions/10105935/how-to-convert-an-int-value-to-string-in-go
// https://stackoverflow.com/questions/47723193/panic-json-cannot-unmarshal-array-into-go-value-of-type-main-structure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	loggly "github.com/jamespearly/loggly"
)

type weatherData struct {
	ID            int     `json:"id"`
	DATE          string  `json:"applicable_date"`
	TEMPERATURE   float64 `json:"the_temp`
	WEATHERSTATUS string  `json:"weather_state_name"`
}

type dbitem struct {
	Id   string
	Time string
	Data []weatherData
}

// AWS.config.update({
// 	endpoint: "https://dynamodb.us-east-1.amazonaws.com"}
//  );

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func main() {

	// os.Setenv("AWS_ACCESS_KEY_ID", "")
	// os.Setenv("AWS_SECRET_ACCESS_KEY", "")
	// os.Setenv("AWS_SESSION_TOKEN", "")

	//https://qvault.io/golang/range-over-ticker-in-go-with-immediate-first-tick/
	ticker := time.NewTicker(time.Minute)
	t := time.Now()

	for ; true; <-ticker.C {
		// currentTime := time.Now()

		var tag string = "firstapplication"

		client := loggly.New(tag)

		// resp, err := http.Get("https://www.metaweather.com/api/location/2459115/2021/9/13/")

		// resp, err := http.NewRequest("GET", "https://www.metaweather.com/api/location/2459115", nil)

		// year := strconv.Itoa(t.Year())
		month := strconv.Itoa(int(t.Month()))
		day := strconv.Itoa(t.Day())

		// params := "year" + url.QueryEscape(year) + "/" +
		// 	"month" + url.QueryEscape(month) + "/" +
		// 	"day" + url.QueryEscape(day) + "/"

		//use the Sprintf() function to format the string without printing and then store it to another variable
		path := fmt.Sprintf("https://www.metaweather.com/api/location/2459115/2021/" + month + "/" + day)
		// resp, err := http.Get(path)

		// params := ("2006/01/02")

		// path := fmt.Sprintf("https://www.metaweather.com/api/location/2459115/", params)

		resp, err := http.Get(path)

		// resp, err := http.Get("https://www.metaweather.com/api/location/2459115/" + currentTime.Format("2021/01/02"))

		if err != nil {
			client.Send("error", "This is an error message:"+err.Error())
			log.Fatal(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			client.Send("error", "This is an error message:"+err.Error())
			log.Fatal(err)
		}

		var x []weatherData
		err = json.Unmarshal(body, &x)

		if err != nil {
			client.Send("error", "This is an error message:"+err.Error())
			log.Fatal(err)
		}

		log.Printf("%+v", x)

		output := strconv.Itoa(int(len(body)))
		//  output2 := resp.ContentLength

		// Valid Send (no error returned)
		err = client.EchoSend("info", "Success! Data size: "+output)
		fmt.Println("err:", err)

		//################################################ DYNAMO DB STARTS BELOW ########################################################################################
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String("us-east-1")},
		)
		if err != nil {
			fmt.Println("Error starting a new session")
		}
		// Create DynamoDB client
		svc := dynamodb.New(sess)

		var dbitem dbitem
		dbitem.Data = x
		dbitem.Time = time.Now().String()
		dbitem.Id = RandomString(10)

		//set the id, time, weatherdb -> the pass the item to marshal tingy

		av, err := dynamodbattribute.MarshalMap(dbitem)
		if err != nil {
			log.Fatalf("Got error marshalling new weather item: %s", err)
		}

		tableName := "csc_482"

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		}

		_, err = svc.PutItem(input)
		if err != nil {
			log.Fatalf("Got error calling PutItem: %s", err)
		}

		fmt.Println("Successfully added '" + dbitem.Id + "' to table " + tableName)

	}

}
