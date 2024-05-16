package main

import (
	"fmt"
	"log"

	"seqquery.de/yenta"
	"seqquery.de/yenta/internal/model"
	"seqquery.de/yenta/internal/util"
)

func main() {
	s, err := yenta.New("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Printf("Can not create yenta service: %s", err)
	}

	ex, err := yenta.NewExchange("log", "topic", false, false, false, true)
	if err != nil {
		log.Printf("Can not declare exchange: %s", err)
	}

	qin, err := yenta.NewQueue("email", false, true, false, false)
	if err != nil {
		log.Printf("Can not declare queue: %s", err)
	}

	qin2, err := yenta.NewQueue("fibonacci", false, true, false, false)
	if err != nil {
		log.Printf("Can not declare queue: %s", err)
	}

	fun := func(arg map[string]interface{}) {
		var n int
		temp := arg["n"]
		if temp == nil {
			fmt.Println("error provided data does not contain n")
			return
		}
		n = int(temp.(float64))
		res, err := util.Fibonacci(n)
		if err != nil {
			log.Printf("error calculating the results %s", err)
		}
		log.Printf("Calculated the results: %+v\n", res)
	}

	s.Subscribe(*ex, qin.(model.Queue), "ubot", fun)
	s.Subscribe(*ex, qin2.(model.Queue), "ubot", fun)
	err = s.Run()
	fmt.Println(err)
}

/*
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"seqquery.de/worker"
	"seqquery.de/worker/mill"
	"seqquery.de/yenta"
)

var env = yenta.New()

func (env *yenta.Yenta) Receive(path string) error {
	conn, err := amqp.Dial(path)
	if err != nil {
		log.Fatalf("Connection to RabbitMQ failed: %s", err)
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to create channel: %s", err)
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %s", err)
		return err
	}

	msg, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to consume message: %s", err)
		return err
	}

	for m := range msg {
		//log.Printf("[<-] Received message: '%s'\n", m.Body)
		job := worker.NewJob()
		arg := make(map[string]interface{})
		json.Unmarshal(m.Body, &arg)
		job.Args = arg
		env.jobs <- job
	}
	defer ch.Close()
	return err
}

func (env *Envoy) Send(path string) {
	con, err := amqp.Dial(path)
	if err != nil {
		log.Fatal("Connection error")
	}
	ch, err := con.Channel()
	if err != nil {
		log.Printf("Can not open the channel: %s", err)
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Can not declare the queue: %s", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var forever chan struct{}

	go func() {
		for {
			body := make(map[string]interface{})
			body["message"] = "Hi!"
			body["time"] = time.Now()
			args, _ := json.Marshal(body)
			err = ch.PublishWithContext(
				ctx,
				"",
				q.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        args,
				},
			)
			if err != nil {
				log.Printf("Failed to publish message: %s", err)
				return
			}
			log.Println("Message sent!")
			time.Sleep(5000 * time.Millisecond)

		}
	}()
	<-forever
}
*/
