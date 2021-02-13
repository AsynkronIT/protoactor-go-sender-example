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
		// ctx.Send() do not set sender information and hence the recipient can not respond.
		ctx.Request(p.pongPid, &ping{})

	case *pong:
		log.Print("Received pong message")

	}
}

func main() {
	system := actor.NewActorSystem()

	pongProps := actor.PropsFromFunc(func(ctx actor.Context) {
		switch ctx.Message().(type) {
		case *ping:
			log.Print("Received ping message")
			ctx.Respond(&pong{})

		default:

		}
	})
	pongPid := system.Root.Spawn(pongProps)

	pingProps := actor.PropsFromProducer(func() actor.Actor {
		return &pingActor{
			pongPid: pongPid,
		}
	})
	pingPid := system.Root.Spawn(pingProps)

	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt)
	signal.Notify(finish, syscall.SIGTERM)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			system.Root.Send(pingPid, struct{}{})

		case <-finish:
			log.Print("Finish")
			return

		}
	}
}
