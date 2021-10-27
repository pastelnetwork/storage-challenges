package main

import (
	"log"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pastelnetwork/storage-challenges/application/dto"
	"github.com/pastelnetwork/storage-challenges/external/message"
)

func main() {
	clients := actor.NewPIDSet(actor.NewPID("127.0.0.1:8000", "storage-challenge"))
	msg := &dto.StorageChallengeRequest{}
	callback := func(context actor.Context, message interface{}) {
		log.Printf("CALLBACK: recieved response %#v", message)
	}

	c := message.Do(clients, msg, callback)
	defer c()

	console.ReadLine()
}
