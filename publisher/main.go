package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/nats-io/stan.go"
	"github.com/oklog/ulid"
)

func main() {
	natsURL := os.Getenv("NATS_URL")
	clusterID := os.Getenv("NATS_CLUSTER_ID")
	topic := os.Getenv("NATS_TOPIC")
	addr := ":" + os.Getenv("PORT")

	err := startPublisher(natsURL, clusterID, addr, topic)
	if err != nil {
		log.Fatal(err)
	}
}

func startPublisher(natsURL, clusterID, addr, topic string) error {
	publisher, err := nats.NewStreamingPublisher(
		nats.StreamingPublisherConfig{
			ClusterID: clusterID,
			ClientID:  "publisher",
			StanOptions: []stan.Option{
				stan.NatsURL(natsURL),
			},
			Marshaler: nats.GobMarshaler{},
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		return err
	}

	h := handler{topic, publisher}
	http.HandleFunc("/", h.Handle)

	log.Print("Listening on ", addr)
	return http.ListenAndServe(addr, nil)
}

type handler struct {
	topic     string
	publisher message.Publisher
}

func (h handler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	err := h.publish(w, r)
	if err != nil {
		log.Println("Error: ", err)
		w.WriteHeader(500)
		return
	}
}

func (h handler) publish(w http.ResponseWriter, r *http.Request) error {
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	uuid := ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader)
	msg := message.NewMessage(uuid.String(), payload)

	if err := h.publisher.Publish(h.topic, msg); err != nil {
		return err
	}

	_, err = fmt.Fprint(w, "Sent message: ", string(payload), " with ID ", msg.UUID, "\n")
	if err != nil {
		return err
	}

	return nil
}
