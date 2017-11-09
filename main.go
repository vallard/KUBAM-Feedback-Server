/* KUBAM Feedback.  Microservice that runs that handles post requests and then forwards to a spark room */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/joeshaw/envdecode"
	"github.com/vallard/spark"
)

type Feedback struct {
	Message string `json:"message"`
}

var ev struct {
	SparkRoom  string `env:"SPARK_ROOM,required"`
	SparkToken string `env:"SPARK_TOKEN,required"`
}

var s *spark.Spark

func handleFeedback(s *spark.Spark, fb Feedback, roomId string) error {
	log.Printf("Message: %s", fb.Message)
	newMessage := spark.Message{
		RoomId: roomId,
		Text:   fb.Message,
	}
	m, err := s.CreateMessage(newMessage)
	if err != nil {
		log.Printf("Unable to create message.\nM: %v\n", m)
	}
	return err
}

func main() {

	if err := envdecode.Decode(&ev); err != nil {
		log.Fatalf("Environment Decode Problem: %v\n", err)
	}
	s = spark.New(ev.SparkToken)
	http.HandleFunc("/v1/feedback", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got a request:\n  %v\n\n", r)
		if r.Method == "POST" {
			decoder := json.NewDecoder(r.Body)
			for {
				var fb Feedback
				if err := decoder.Decode(&fb); err != nil {
					log.Print(err)
					break
				}
				// do something with the message.
				//handleFeedback(s, fb, ev.SparkRoom)
				handleFeedback(s, fb, "Y2lzY29zcGFyazovL3VzL1JPT00vYzJjNDY3MDAtYzhhMS0xMWU2LThmNmEtZTlmZTYyZjkwMzU1")

			}
		}
		if r.Method == "GET" {
			fmt.Fprintf(w, "I'm alive and waiting for feedback")

		}
	})

	log.Print("Kubam Feedback is Listening on port 9999")
	log.Print("call me with: curl -X POST -d '{\"message\" : \"Kubam is awesome.\" }' localhost:9999/v1/feedback ")
	http.ListenAndServe(":9999", nil)
}
