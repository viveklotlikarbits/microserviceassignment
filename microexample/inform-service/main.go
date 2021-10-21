package main

import (
	"fmt"
	"log"
	"github.com/streadway/amqp"
	"io/ioutil"
	"net/http"
	"encoding/json"	
	"bytes"
	"strings"
)

func failOnError(err error, msg string) {
	if err != nil {
		// Exit the program.
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	// 'rabbitmq-server' is the network reference we have to the broker, 
	// thanks to Docker Compose.
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq-server:5672/")
	failOnError(err, "Error connecting to the broker")
	// Make sure we close the connection whenever the program is about to exit.
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	// Make sure we close the channel whenever the program is about to exit.
	defer ch.Close()
	
	exchangeName := "product_order"
	bindingKey   := "product.order.*"

	// Create the exchange if it doesn't already exist.
	err = ch.ExchangeDeclare(
			exchangeName, 	// name
			"topic",  		// type
			true,         	// durable
			false,
			false,
			false,
			nil,
	)
	failOnError(err, "Error creating the exchange")
	
	// Create the queue if it doesn't already exist.
	// This does not need to be done in the publisher because the
	// queue is only relevant to the consumer, which subscribes to it.
	// Like the exchange, let's make it durable (saved to disk) too.
	q, err := ch.QueueDeclare(
			"",    // name - empty means a random, unique name will be assigned
			true,  // durable
			false, // delete when unused
			false, 
			false, 
			nil,   
	)
	failOnError(err, "Error creating the queue")

	// Bind the queue to the exchange based on a string pattern (binding key).
	err = ch.QueueBind(
			q.Name,       // queue name
			bindingKey,   // binding key
			exchangeName, // exchange
			false,
			nil,
	)
	failOnError(err, "Error binding the queue")

	// Subscribe to the queue.
	msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer id - empty means a random, unique id will be assigned
			false,  // auto acknowledgement of message delivery
			false,  
			false,  
			false,  
			nil,
	)
	failOnError(err, "Failed to register as a consumer")
	//app := "curl"
	//arg0 := "-k -H "Content-Type: application/json\" -H \"x-apikey: e2ef39ac6fc3ac51270e07553ce2ba3f9b83f\" -X POST -d \'{\"subject":\""
	//arg1 = "test2"
	//arg2 := "\"}\' \'https://studentsubject-d4fe.restdb.io/rest/subject\'"
	type Subject struct {
	// defining struct variables
		name      string
	}	


	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received Enrollment for : %s", d.Body)

			client := &http.Client{}
			var subj Subject
			//rawIn:= json.RawMessage(d.Body)
			//bytes1, err:= rawIn.MarshalJSON()
			//if err != nil {
				// if error is not nil
				// print error
				//fmt.Println(err)
			//}
			s := strings.Split(string(d.Body), "\"")
			log.Printf("Received Enrollment1 for : %s", s)
			log.Printf("Received Enrollment2 for : %s", s[3])

			//err1 := json.Unmarshal(bytes1, &subj)
			err1 := json.Unmarshal(d.Body, &subj)
			if err1 != nil {
				// if error is not nil
				// print error
				fmt.Println(err1)
			}
			log.Printf("Received Enrollment for : %s", subj)

			var jsonstr = "{\"subject\":\"" + s[3] + "\"}"
			var jsonData = []byte(jsonstr) 

			log.Printf("Received Enrollment for : %s", jsonData)

			req, err := http.NewRequest("POST", "https://studentsubject-d4fe.restdb.io/rest/subject", bytes.NewBuffer(jsonData))
			if err != nil {
				// handle err
				fmt.Println("Error in Request Creation")
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Apikey", "e2ef39ac6fc3ac51270e07553ce2ba3f9b83f")

			response, err := client.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error in connecting client")
			}
			defer response.Body.Close()

			fmt.Println("response Status:", response.Status)
			fmt.Println("response Headers:", response.Header)
			body, _ := ioutil.ReadAll(response.Body)
			fmt.Println("response Body:", string(body))

			// Update the user data on the service's 
			// associated datastore using a local transaction...

			// The 'false' indicates the success of a single delivery, 'true' would mean that
			// this delivery and all prior unacknowledged deliveries on this channel will be
			// acknowledged, which I find no reason for in this example.
			d.Ack(false)
		}
	}()
	
	fmt.Println("Inform Service listening for Enrollments...")
	
	// Block until 'forever' receives a value, which will never happen.
	<-forever
}
