package main

import (
	"log"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	appActor "github.com/pastelnetwork/storage-challenges/application/actor"
	"github.com/pastelnetwork/storage-challenges/domain/service"
	"github.com/pastelnetwork/storage-challenges/external/message"
	"github.com/pastelnetwork/storage-challenges/external/storage"
)

func main() {
	remoter := message.NewRemoter(actor.NewActorSystem(), message.Config{Host: "localhost", Port: 9003})
	dommainService := service.NewStorageChallenge(service.Config{Remoter: remoter})
	store, err := storage.NewStore(storage.Config{})
	if err != nil {
		log.Fatal("storage.NewStore", err)
	}
	a := appActor.NewStorageChallengeActor(dommainService, store)
	_, err = remoter.RegisterActor(a, "storage-challenge")
	if err != nil {
		log.Fatal("remoter.RegisterActor", err)
	}
	remoter.Start()
	defer remoter.GracefulStop()
	console.ReadLine()
}
