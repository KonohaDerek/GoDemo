package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/philippseith/signalr"
	"github.com/rs/zerolog/log"
)

type chat struct {
	signalr.Hub
}

func (c *chat) OnConnected(connectionID string) {
	fmt.Printf("%s connected\n", connectionID)
	c.Groups().AddToGroup("group", connectionID)
}

func (c *chat) OnDisconnected(connectionID string) {
	fmt.Printf("%s disconnected\n", connectionID)
	c.Groups().RemoveFromGroup("group", connectionID)
}

func (c *chat) Broadcast(message string) {
	// Broadcast to all clients
	c.Clients().Group("group").Send("receive", message)
}

func (c *chat) Echo(message string) {
	c.Clients().Caller().Send("receive", message)
}

func (c *chat) Panic() {
	panic("Don't panic!")
}

func (c *chat) RequestAsync(message string) <-chan map[string]string {
	r := make(chan map[string]string)
	go func() {
		defer close(r)
		time.Sleep(4 * time.Second)
		m := make(map[string]string)
		m["ToUpper"] = strings.ToUpper(message)
		m["ToLower"] = strings.ToLower(message)
		m["len"] = fmt.Sprint(len(message))
		r <- m
	}()
	return r
}

func (c *chat) RequestTuple(message string) (string, string, int) {
	return strings.ToUpper(message), strings.ToLower(message), len(message)
}

func (c *chat) DateStream() <-chan string {
	r := make(chan string)
	go func() {
		defer close(r)
		for i := 0; i < 50; i++ {
			r <- fmt.Sprint(time.Now().Clock())
			time.Sleep(time.Second)
		}
	}()
	return r
}

func (c *chat) UploadStream(upload1 <-chan int, factor float64, upload2 <-chan float64) {
	ok1 := true
	ok2 := true
	u1 := 0
	u2 := 0.0
	c.Echo(fmt.Sprintf("f: %v", factor))
	for {
		select {
		case u1, ok1 = <-upload1:
			if ok1 {
				c.Echo(fmt.Sprintf("u1: %v", u1))
			} else if !ok2 {
				c.Echo("Finished")
				return
			}
		case u2, ok2 = <-upload2:
			if ok2 {
				c.Echo(fmt.Sprintf("u2: %v", u2))
			} else if !ok1 {
				c.Echo("Finished")
				return
			}
		}
	}
}

func (c *chat) Abort() {
	fmt.Println("Abort")
	c.Hub.Abort()
}

func runHTTPServer(address string, hub signalr.HubInterface) {
	server, _ := signalr.NewServer(context.TODO(), signalr.SimpleHubFactory(hub),
		// signalr.Logger(zerolog.New(os.Stdout).With().Timestamp().Logger(), false),
		signalr.KeepAliveInterval(2*time.Second))
	router := http.NewServeMux()
	server.MapHTTP(signalr.WithHTTPServeMux(router), "/chat")

	fmt.Printf("Serving public content from the embedded filesystem\n")
	// router.Handle("/", http.FileServer(http.FS(public.FS)))
	fmt.Printf("Listening for websocket connections on http://%s\n", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatal().
			Err(err).
			Msg("invalid ListenAndServe")
	}
	// if err := http.ListenAndServe(address, middleware.LogRequests(router)); err != nil {
	// 	log.Fatal().
	// 		Err(err).
	// 		Msg("invalid ListenAndServe")
	// }
}

func runHTTPClient(address string, receiver interface{}) error {
	c, err := signalr.NewClient(context.Background(), nil,
		signalr.WithReceiver(receiver),
		signalr.WithConnector(func() (signalr.Connection, error) {
			creationCtx, _ := context.WithTimeout(context.Background(), 2*time.Second)
			return signalr.NewHTTPConnection(creationCtx, address)
		}),
		// signalr.Logger(kitlog.NewLogfmtLogger(os.Stdout), false)
	)
	if err != nil {
		return err
	}
	c.Start()
	fmt.Println("Client started")
	return nil
}

type receiver struct {
	signalr.Receiver
}

func (r *receiver) Receive(msg string) {
	fmt.Println(msg)
	// The silly client urges the server to end his connection after 10 seconds
	r.Server().Send("abort")
}

func main() {
	hub := &chat{}

	//go runTCPServer("127.0.0.1:8007", hub)
	go runHTTPServer("localhost:8086", hub)
	<-time.After(time.Millisecond * 2)
	go func() {
		fmt.Println(runHTTPClient("http://localhost:8086/chat", &receiver{}))
	}()
	ch := make(chan struct{})
	<-ch
}
