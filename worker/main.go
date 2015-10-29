package main

import (
	"github.com/garyburd/redigo/redis"
	"time"

	"encoding/json"
	"log"
	"fmt"
	"net/http"
)

var redisChan chan Candy

type Candy struct {
	Name, Object string
	Price        float64
	Time         time.Time
}

func redisSubscriber() {
	conn, err := redis.Dial("tcp", "redis:6379")
	if err != nil {
		panic("I don't want to live on this planet anymore")
	}
	psc := redis.PubSubConn{conn}
	psc.Subscribe("candy")
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
			if v.Channel == "candy" {
				var c Candy
				err := json.Unmarshal(v.Data, c)
				if (err != nil) {
					log.Printf("Seems our redis is sick! In the evening we'll get some schnaps to ease the pain!")
				}
				redisChan <- c
			}
		case redis.Subscription:
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			log.Println(v)
		}
	}
}

func elasticUpdater() {


	// Add a document to the index

	for c := range redisChan {
		buf, err := json.Marshal(c)
		if err != nil {
			log.Print(err)
			continue
		}
		resp, err := http.Post("http://elasticsearch:9200/candyshop/candy/", "application/json", &buf)
		if err != nil {
			log.Println(err)
		}
		if resp != nil {
			log.Println(resp)
		}
	}
}

func main() {
	time.Sleep(time.Second * 20)
	redisChan = make(chan Candy, 10)
	go redisSubscriber()
	elasticUpdater()


}
