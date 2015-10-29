package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
	"gopkg.in/olivere/elastic.v3"

	"encoding/json"
	"log"
	"fmt"
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
	client, err := elastic.NewClient(elastic.SetURL("http://elasticsearch:9200"))
	if err != nil {
		// Handle error
	}

	_, err = client.CreateIndex("twitter").Do()
	if err != nil {
		// Handle error
		panic(err)
	}

	// Add a document to the index

	for c := range redisChan {
		_, err = client.Index().
		Index("candyshop").
		Type("candy").
		BodyJson(c).
		Do()
		if err != nil {
			// Handle error
			log.Println(err)
		}
	}
}

func main() {
	redisChan = make(chan Candy, 10)
	go redisSubscriber()
	elasticUpdater()


}
