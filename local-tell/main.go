package main

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type pong struct {
}

type ping struct {
}

type pingActor struct {
	pongPid *actor.PID
}

func (p *pingActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case struct{}:
		// Below both do not set ctx.Self() as sender,
		// and hence the recipient has no knowledge of the sender
		// even though the message is sent from another actor.
		//
		//p.pongPid.Tell(&ping{})
		ctx.Tell(p.pongPid, &ping{})

	case *pong:
		log.Print("Received pong message")

	}
}

func main() {
	pongProps := actor.FromFunc(func(ctx actor.Context) {
		switch ctx.Message().(type) {
		case *ping:
			log.Print("Received ping message")
			// Below both fail, but their behavior slightly differ.
			// ctx.Sender().Tell() panics and recovers because this tries to call nil sender's Tell method;
			// while ctx.Respond() checks the presence of sender and redirects the message to dead letter process
			// when sender is absent.
			//
			ctx.Sender().Tell(&pong{})
			//ctx.Respond(&pong{})

		default:

		}
	})
	pongPid := actor.Spawn(pongProps)

	pingProps := actor.FromProducer(func() actor.Actor {
		return &pingActor{
			pongPid: pongPid,
		}
	})
	pingPid := actor.Spawn(pingProps)

	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt)
	signal.Notify(finish, syscall.SIGTERM)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pingPid.Tell(struct{}{})

		case <-finish:
			return
			log.Print("Finish")

		}
	}
}