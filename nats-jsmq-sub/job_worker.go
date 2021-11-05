package main

import (
	"fmt"
	"log"
	"runtime"
	"time"
	"github.com/nats-io/nats.go"
)

func main() {
	var urls = nats.DefaultURL

	// Connect Options.
	opts := []nats.Option{nats.Name("NATS JetStream Subscriber")}
	opts = setupConnOptions(opts)

	// Connect to NATS
	nc, err := nats.Connect(urls, opts...)
	if err != nil {
		log.Fatal(err)
	}

	// Use the JetStream context to produce and consumer messages
	// that have been persisted.
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Waiting for job...");
	
	i := 0
	js.QueueSubscribe("jobs", "group", func(msg *nats.Msg) {
		i++
		log.Printf("[#%d] Start Working on... [%s]: '%s'", i, msg.Subject, string(msg.Data))		
		<-time.After(5*time.Second)
		msg.Ack()
		log.Printf("[#%d] Done  Working on... [%s]: '%s'", i, msg.Subject, string(msg.Data))		
	}, nats.ManualAck())

	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}
	
	runtime.Goexit()
}

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	return opts
}
