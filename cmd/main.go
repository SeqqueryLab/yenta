package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/utubun/yenta"
)

func main() {
	s, err := yenta.New("amqp://guest:guest@localhost:5672")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("New service is created")
	b, _ := os.ReadFile("internal/assets/config.json")
	err = s.Configure(b)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Service: %+v\n", s)

	go func() {
		for {
			body, _ := json.Marshal(struct {
				Time time.Time
				Text string
			}{
				Time: time.Now(),
				Text: "Hi there",
			})
			err = s.Post("human numbers", []string{"hot e", "sexy e", "I rot"}, false, false, "text/plain", body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Sent message: %s\n", string(body))
			time.Sleep(5 * time.Second)
		}
	}()

	err = s.Get("main", "e-queue", true, false, false, false, nil, logger)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	fmt.Print(err)
}

type body struct {
	Time time.Time
	Text string
}

func logger(m []byte) any {
	var res body
	json.Unmarshal(m, &res)
	log.Printf("Logging the results %+v", res)
	return nil
}
