package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/SeqqueryLab/yenta"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	s := yenta.NewService("amqp://guest:guest@localhost:5672/")
	log.Println("New service is created")

	//logs := yenta.NewExchange("logs", "fanout", true, false, false, false, nil)
	nums := yenta.NewExchange("numbers", "topic", true, false, false, true, nil)

	//qsq := yenta.NewQueue("squares", true, false, false, false, nil)

	qfb := yenta.NewQueue("fibonacci", true, false, false, false, nil)

	s.Worker(nums, qfb, []string{"log", "fibonacci"}, logger)

	err := s.Run()
	fmt.Println(err)
}

func logger(arg interface{}) interface{} {
	body := arg.(amqp091.Delivery).Body
	res := make(map[string]interface{})

	json.Unmarshal(body, &res)
	//fibbo, err := util.Fibonacci(int(res["n"].(float64)))
	//if err != nil {
	//	log.Printf("Error: %s", err)
	//}
	log.Printf("Logging the results %+v", res)
	return nil
}
