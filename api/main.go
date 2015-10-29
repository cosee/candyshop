package main

import (
    "net/http"
    "github.com/garyburd/redigo/redis"
	"time"

	"io"
	"encoding/json"
	"log"
)

var redisChan chan Candy

type Candy struct {
	Name, Object string
	Price float64
	Time time.Time
}

func redisPublisher(){
	conn, err := redis.Dial("tcp", "redis:6379")
	if err != nil {
		panic("I don't want to live on this planet anymore")
	}
	for c := range redisChan{
		buf, err := json.Marshal(c)
		if err != nil {
			log.Print(err)
			continue
		}
		conn.Do("PUBLISH", "candy", buf)
	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "Go away! Watch the tutorial at https://www.youtube.com/watch?v=oHg5SJYRHA0\n")
		return
	}
	var body []byte
	var candy Candy
    io.ReadFull(r.Body, body)
	defer r.Body.Close()
	err := json.Unmarshal(body, candy)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Could not understand you\n")
		return
	}
	candy.Time = time.Now()

	//sanity checks go here
	redisChan <- candy

}

func main() {
	redisChan = make(chan Candy, 10)
	go redisPublisher()

    http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
