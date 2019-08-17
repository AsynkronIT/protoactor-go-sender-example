package main

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/cluster/consul"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/oklahomer/protoactor-go-sender-example/cluster/messages"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var cnt uint64 = 0

type pingActor struct {
	cnt uint
}

func (p *pingActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case struct{}:
		cnt += 1
		ping := &messages.Ping{
			Cnt: cnt,
		}

		grainPid, statusCode := cluster.Get("ponger-1", "Ponger")
		if statusCode != remote.ResponseStatusCodeOK {
			log.Printf("Get PID failed with StatusCode: %v", statusCode)
			return
		}
		// Below do not set ctx.Self() as sender,
		// and hence the recipient has no knowledge of the sender
		// even though the message is sent from another actor.
		//
		// ctx.Send(grainPid, ping)
		ctx.Request(grainPid, ping)

	case *messages.Pong:
		// Never comes here.
		// The recipient can not refer to the sender.
		// Instead the cluster grain leaves logs as below:
		// YYYY/MM/DD hh:mm:ss Received ping message
		// YYYY/MM/DD hh:mm:ss [ACTOR] [DeadLetter] pid="nil" message="&Pong{Cnt:2,}" sender="nil"
		//
		log.Print("Received pong message")

	}
}

func main() {
	cp, err := consul.New()
	if err != nil {
		log.Fatal(err)
	}
	cluster.Start("cluster-example", "127.0.0.1:8081", cp)

	rootContext := actor.EmptyRootContext

	pingProps := actor.PropsFromProducer(func() actor.Actor {
		return &pingActor{}
	})
	pingPid := rootContext.Spawn(pingProps)

	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt)
	signal.Notify(finish, syscall.SIGTERM)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rootContext.Send(pingPid, struct{}{})

		case <-finish:
			log.Print("Finish")
			return

		}
	}
}
