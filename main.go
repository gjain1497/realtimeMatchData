package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gopkg.in/antage/eventsource.v1"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type MatchData struct {
	Timestamp int64
	Data      struct {
		ID        int64  `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Avatar    string `json:"avatar"`
	} `json:"data"`
	Support struct {
		URL  string `json:"url"`
		Text string `json:"text"`
	} `json:"support"`
}

var mydata MatchData

var rdb *redis.Client
var ctx = context.Background()

func initClient() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
}

func setKey(key string, value string) {
	err := rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}

func getKey(key string) string {
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)
	if err == redis.Nil {
		fmt.Println("key does not exist")
	}
	return val
}

func receiveData() {
	resp, err := http.Get("https://reqres.in/api/users/2")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	err1 := json.Unmarshal(body, &mydata)
	if err1 != nil {
		log.Fatalln(err1)
	}
	mydata.Timestamp = time.Now().Unix()
	fmt.Println("data received : ", mydata)
	setKey("abcd", string(body))
}

func getData() {
	es := eventsource.New(nil, nil) //keeping default settings and custom headers
	defer es.Close()

	mydataStr := getKey("abcd")
	fmt.Println("mydataStr ", mydataStr)

	http.Handle("/events", es)

	go func() {
		id := 1
		for {
			es.SendEventMessage(mydataStr, "tick-event", strconv.Itoa(id))
			id++
			time.Sleep(2 * time.Second)
		}
	}()

	log.Fatal(http.ListenAndServe(":8080", es))
}

func main() {
	initClient()
	runCronJon()
	getData()
}
