package main

import (
	"context"
	"fmt"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pastelnetwork/gonode/common/net/credentials"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/storage-challenges/application/dto"
	"github.com/pastelnetwork/storage-challenges/config"
	"github.com/pastelnetwork/storage-challenges/external/message"
	testnodes "github.com/pastelnetwork/storage-challenges/test_nodes"
	appcontext "github.com/pastelnetwork/storage-challenges/utils/context"
)

func main() {
	cfg := config.Config{}
	if err := cfg.Load(); err != nil {
		panic(fmt.Sprintf("could not load config data: %v", err))
	}
	if cfg.Remoter == nil {
		cfg.Remoter = &message.Config{}
	}
	pastelClient := testnodes.NewMockPastelClient()
	remoter := message.NewRemoter(actor.NewActorSystem(), *cfg.Remoter.
		WithClientSecureCreds(credentials.NewClientCreds(pastelClient, &alts.SecInfo{PastelID: "mock pastel id", PassPhrase: "mock passphrase", Algorithm: "mock algorithm"})).
		WithServerSecureCreds(credentials.NewServerCreds(pastelClient, &alts.SecInfo{PastelID: "mock pastel id", PassPhrase: "mock passphrase", Algorithm: "mock algorithm"})))
	remoter.Start()
	defer remoter.GracefulStop()
	pid := actor.NewPID("localhost:9000", "storage-challenge")
	remoter.Send(appcontext.FromContext(context.Background()), pid, &dto.GenerateStorageChallengeRequest{ChallengingMasternodeId: "jXEZVtIEVmSkYw0v8qGjsBrrELBOPuedNYMctelLWSlw6tiVNljFMpZFir30SN9r645tEAKwEAYfKR3o4Ek5YM", ChallengesPerMasternodePerBlock: 1})
	console.ReadLine()
}
