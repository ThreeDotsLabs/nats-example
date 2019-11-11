package main

import (
	"context"
	"log"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/nats-io/stan.go"
)

func main() {
	natsURL := os.Getenv("NATS_URL")
	clusterID := os.Getenv("NATS_CLUSTER_ID")
	topic := os.Getenv("NATS_TOPIC")

	err := startSubscriber(natsURL, clusterID, topic)
	if err != nil {
		log.Fatal(err)
	}
}

func startSubscriber(natsURL, clusterID, topic string) error {
	logger := watermill.NewStdLogger(false, false)

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return err
	}

	subscriber, err := nats.NewStreamingSubscriber(
		nats.StreamingSubscriberConfig{
			ClusterID:   clusterID,
			ClientID:    "subscriber",
			QueueGroup:  "example-group",
			DurableName: "example-durable",
			StanOptions: []stan.Option{
				stan.NatsURL(natsURL),
			},
			Unmarshaler: nats.GobMarshaler{},
		},
		logger,
	)
	if err != nil {
		return err
	}

	router.AddMiddleware(middleware.Recoverer)

	router.AddNoPublisherHandler(
		"messages_handler",
		topic,
		subscriber,
		handler,
	)

	log.Print("Subscribed for messages")

	return router.Run(context.Background())
}

func handler(msg *message.Message) error {
	log.Printf("received message: %s, payload: %s", msg.UUID, string(msg.Payload))
	return nil
}
