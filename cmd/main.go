package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
	"github.com/utubun/yenta"
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

func square(arg map[string]interface{}) map[string]interface{} {
	var n float64
	temp := arg["n"]
	if temp == nil {
		fmt.Println("error provided data does not contain n")
		return nil
	}
	n = temp.(float64)
	n = n * n

	return map[string]interface{}{"n": n}
}

func logger(arg interface{}) interface{} {
	body := arg.(amqp091.Delivery).Body
	res := make(map[string]interface{})

	json.Unmarshal(body, &res)
	log.Printf("Logging results: %+v", res)
	return nil
}

/* func fibonacci(arg interface{}) interface{} {
	temp := arg["n"]
	if temp == nil {
		fmt.Println("error provided data does not contain n")
		return nil
	}
	n := int(temp.(float64))
	res, err := util.Fibonacci(n)
	if err != nil {
		log.Printf("error calculating the results %s", err)
	}
	log.Printf("Calculated the results: %+v\n", res)

	return map[string]interface{}{"res": res}
} */
