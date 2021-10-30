package main

import (
	"context"
	"log"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	appActor "github.com/pastelnetwork/storage-challenges/application/actor"
	"github.com/pastelnetwork/storage-challenges/application/dto"
	"github.com/pastelnetwork/storage-challenges/domain/service"
	"github.com/pastelnetwork/storage-challenges/external/message"
	"github.com/pastelnetwork/storage-challenges/external/storage"
	appcontext "github.com/pastelnetwork/storage-challenges/utils/context"
)

func main() {
	remoter := message.NewRemoter(actor.NewActorSystem(), message.Config{Remoter: message.Address{Host: "localhost", Port: 9000}})
	dommainService := service.NewStorageChallenge(remoter, nil)
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
	err = remoter.Send(appcontext.FromContext(context.Background()), message.ActorProperties{
		Address: "localhost:9001",
		Name:    "storage-challenge",
		Kind:    "storage-challenge",
	}, &dto.StorageChallengeRequest{Data: &dto.StorageChallengeData{
		ChallengeId:             "test-challenge-id",
		ChallengingMasternodeId: "test-challenging-masternode-id",
		MessageId:               "test-message-id",
		MessageType:             dto.MessageType_STORAGE_CHALLENGE_ISSUANCE_MESSAGE,
		ChallengeFile: &dto.StorageChallengeDataChallengeFile{
			FileHashToChallenge:      "test-file-hash-to-challenge",
			ChallengeSliceStartIndex: 0,
			ChallengeSliceEndIndex:   2,
		},
	}})
	if err != nil {
		log.Fatal(err)
	}
	console.ReadLine()
}
