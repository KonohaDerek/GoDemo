package pubsub

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/nats-io/stan.go"
)

func NewNATSPublisher(logger watermill.LoggerAdapter, clusterID, natsURL string) (message.Publisher, error) {
	return nats.NewStreamingPublisher(
		nats.StreamingPublisherConfig{
			ClusterID: clusterID,
			ClientID:  watermill.NewShortUUID(),
			StanOptions: []stan.Option{
				stan.NatsURL(natsURL),
			},
			Marshaler: nats.GobMarshaler{},
		},
		logger,
	)
}
func NewNATSSubscriber(logger watermill.LoggerAdapter, clusterID, clientID, natsURL string) (message.Subscriber, error) {
	return nats.NewStreamingSubscriber(
		nats.StreamingSubscriberConfig{
			ClusterID: clusterID,
			ClientID:  clientID,
			StanOptions: []stan.Option{
				stan.NatsURL(natsURL),
			},
			Unmarshaler: nats.GobMarshaler{},
		},
		logger,
	)
}
