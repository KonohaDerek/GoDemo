package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/spf13/cobra"
	"k-derek.dev/demo/cqrs/internal/cqrs/commands"
	"k-derek.dev/demo/cqrs/internal/cqrs/events"
	"k-derek.dev/demo/cqrs/internal/cqrs/handlers"
	"k-derek.dev/demo/cqrs/internal/pubsub"
)

// ServerCmd 是此程式的Service入口點
var ServerCmd = &cobra.Command{
	Run: run,
	Use: "server",
}

var (
	watermillPubsubType = os.Getenv("WATERMILL_PUBSUB_TYPE")
	natsClusterID       = os.Getenv("NATS_CLUSTER_ID")
	natsURL             = os.Getenv("NATS_URL")
	// You probably want to ship your own implementation of `watermill.LoggerAdapter`.
	logger = watermill.NewStdLogger(false, false) // debug=false, trace=false
	// incomingTopic = "incoming_topic"
	// outgoingTopic = "outgoing_topic"
)

func run(_ *cobra.Command, _ []string) {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		log.Fatal(err)
	}
	// SignalsHandler gracefully shutdowns Router when receiving SIGTERM
	router.AddPlugin(plugin.SignalsHandler)

	// Router level middleware are executed for every message sent to the router
	router.AddMiddleware(
		// CorrelationID will copy the correlation id from the incoming message's metadata to the produced messages
		middleware.CorrelationID,
		// Timeout makes the handler cancel the incoming message's context after a specified time
		middleware.Timeout(time.Second*10),
		// Throttle provides a middleware that limits the amount of messages processed per unit of time
		middleware.NewThrottle(10, time.Second).Middleware,
		// After MaxRetries, the message is Nacked and it's up to the PubSub to resend it
		middleware.Retry{
			MaxRetries: 5,
			Logger:     logger,
		}.Middleware,

		// Recoverer handles panics from handlers
		middleware.Recoverer,
	)

	var publisher message.Publisher
	var subscriber message.Subscriber
	var eventsPublisher message.Publisher
	if watermillPubsubType == "nats" {
		publisher, err = pubsub.NewNATSPublisher(logger, natsClusterID, natsURL)
		if err != nil {
			log.Fatal(err)
		}
		subscriber, err = pubsub.NewNATSSubscriber(logger, natsClusterID, watermill.NewShortUUID(), natsURL)
		if err != nil {
			log.Fatal(err)
		}

		eventsPublisher, err = pubsub.NewNATSPublisher(logger, natsClusterID, natsURL)
		if err != nil {
			panic(err)
		}

	} else {
		pubsub := gochannel.NewGoChannel(
			gochannel.Config{},
			logger,
		)
		publisher = pubsub
		subscriber = pubsub
		eventsPublisher = pubsub
	}

	// cqrs.Facade is facade for Command and Event buses and processors.
	// You can use facade, or create buses and processors manually (you can inspire with cqrs.NewFacade)
	cqrsFacade, err := cqrs.NewFacade(cqrs.FacadeConfig{
		GenerateCommandsTopic: func(commandName string) string {
			// we are using queue RabbitMQ config, so we need to have topic per command type
			return commandName
		},
		CommandHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.CommandHandler {
			return []cqrs.CommandHandler{
				handlers.BookRoomHandler{EventBus: eb},
				handlers.OrderBeerHandler{EventBus: eb},
			}
		},
		CommandsPublisher: publisher,
		CommandsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
			log.Printf("[CommandsSubscriberConstructor] %s", handlerName)
			// we can reuse subscriber, because all commands have separated topics
			return subscriber, nil
		},
		GenerateEventsTopic: func(eventName string) string {
			// because we are using PubSub RabbitMQ config, we can use one topic for all events
			return "events"
			// we can also use topic per event type
			// return eventName
		},
		EventHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.EventHandler {
			return []cqrs.EventHandler{
				events.OrderBeerOnRoomBooked{CommandBus: cb},
				// NewBookingsFinancialReport(),
			}
		},
		EventsPublisher: eventsPublisher,
		EventsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
			log.Printf("[EventsSubscriberConstructor]")
			return subscriber, nil
		},
		Router:                router,
		CommandEventMarshaler: cqrs.JSONMarshaler{},
		Logger:                logger,
	})
	if err != nil {
		panic(err)
	}

	// Producing some incoming messages in background
	// publish BookRoom commands every second to simulate incoming traffic
	go publishCommands(cqrsFacade.CommandBus())

	// Run the router
	ctx := context.Background()
	if err := router.Run(ctx); err != nil {
		log.Fatal(err)
	}
}

func publishCommands(commandBus *cqrs.CommandBus) func() {
	i := 0
	for {
		i++

		startDate := time.Now()

		endDate := time.Now().Add(time.Hour * 24 * 3)

		bookRoomCmd := &commands.BookRoomCmd{
			RoomId:    fmt.Sprintf("%d", i),
			GuestName: "John",
			StartDate: startDate,
			EndDate:   endDate,
		}
		if err := commandBus.Send(context.Background(), bookRoomCmd); err != nil {
			panic(err)
		}

		time.Sleep(time.Second)
	}
}
