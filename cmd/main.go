package main

import (
	"fmt"
	"log"
	"time"

	"github.com/utubun/yenta"
	"github.com/utubun/yenta/internal/util"
)

func main() {
	s, err := yenta.New("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Printf("Can not create yenta service: %s", err)
	}

	logs, err := yenta.NewExchange("logs", "topic", false, false, false, true)
	if err != nil {
		log.Printf("Can not declare exchange: %s", err)
	}
	nums, err := yenta.NewExchange("numbers", "topic", false, false, false, true)
	if err != nil {
		log.Printf("Can not declare exchange: %s", err)
	}

	qsq, err := yenta.NewQueue("squares", false, false, false, false)
	if err != nil {
		log.Printf("Can not declare queue: %s", err)
	}

	qfb, err := yenta.NewQueue("fibonacci", false, false, false, false)
	if err != nil {
		log.Printf("Can not declare queue: %s", err)
	}

	quo, err := yenta.NewQueue("log", false, true, false, false)
	if err != nil {
		log.Printf("Can not create queue: %s", err)
	}

	cfg := yenta.Config{
		{
			Consumer: yenta.Consumer{Exchange: nums, Queue: qfb, Rout: "fibonacci", Worker: fibonacci},
			Producer: yenta.Producer{Exchange: logs, Queue: quo, Rout: "", Worker: logger},
		},
		{
			Consumer: yenta.Consumer{Exchange: nums, Queue: qsq, Rout: "squares", Worker: square},
			Producer: yenta.Producer{Exchange: logs, Queue: quo, Rout: "", Worker: logger},
		},
	}

	s.Add(cfg)

	fmt.Printf("The service: %+v\n", s)
	err = s.Run()
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

func logger(arg map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})

	res["time"] = time.Now()
	res["message"] = "results are ready"
	res["data"] = arg

	return res
}

func fibonacci(arg map[string]interface{}) map[string]interface{} {
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
}
