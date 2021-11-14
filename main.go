package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pastelnetwork/gonode/common/net/credentials"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/pastel"
	appactor "github.com/pastelnetwork/storage-challenges/application/actor"
	"github.com/pastelnetwork/storage-challenges/config"
	"github.com/pastelnetwork/storage-challenges/domain/service"
	"github.com/pastelnetwork/storage-challenges/external/message"
	"github.com/pastelnetwork/storage-challenges/external/repository"
	"github.com/pastelnetwork/storage-challenges/external/storage"
	testnodes "github.com/pastelnetwork/storage-challenges/test_nodes"
)

func main() {
	var migrate, seed, test bool
	flag.Bool("migrate", false, "migration only")
	flag.Bool("migrate-seed", false, "migration with seeding dummy data")
	flag.Bool("test", false, "run node in test mode for debugging purpose")
	flag.Parse()
	flag.VisitAll(func(f *flag.Flag) {
		if f.Name == "migrate" {
			if f.Value.String() == "true" {
				migrate = true
			}
		}
		if f.Name == "migrate-seed" {
			if f.Value.String() == "true" {
				migrate = true
				seed = true
			}
		}
		if f.Name == "test" {
			if f.Value.String() == "true" {
				test = true
			}
		}
	})

	cfg := config.Config{}
	if err := cfg.Load(); err != nil {
		panic(fmt.Sprintf("could not load config data: %v", err))
	}
	if cfg.Database == nil {
		panic("database configuration not found")
	}
	store, err := storage.NewStore(*cfg.Database)
	if err != nil {
		panic(fmt.Sprintf("could not connect to database %v", err))
	}

	if migrate {
		testnodes.AutoMigrate(store.GetDB(), seed)
		return
	}

	if cfg.Remoter == nil {
		cfg.Remoter = &message.Config{}
	}
	if cfg.PastelClient == nil {
		cfg.PastelClient = pastel.NewConfig()
	}

	var pastelClient pastel.Client
	var secInfo *alts.SecInfo
	if test {
		log.Println("USING TEST PASTEL CLIENT")
		pastelClient = testnodes.NewMockPastelClient(store.GetDB())
		secInfo = &alts.SecInfo{
			PastelID:   "mock pastel id",
			PassPhrase: "mock mock passphrase",
			Algorithm:  "mock algorithm",
		}
	} else {
		pastelClient = pastel.NewClient(cfg.PastelClient)
		secInfo = &alts.SecInfo{
			PastelID:   cfg.MasternodePastelID,
			PassPhrase: cfg.MasternodePastelPassphrase,
			Algorithm:  "ed448",
		}
	}

	remoter := message.NewRemoter(
		actor.NewActorSystem(),
		*cfg.Remoter.
			WithClientSecureCreds(credentials.NewClientCreds(pastelClient, secInfo)).
			WithServerSecureCreds(credentials.NewServerCreds(pastelClient, secInfo)),
	)
	remoter.Start()
	defer remoter.GracefulStop()

	repo := repository.New()

	domainService := service.NewStorageChallenge(service.Config{
		Remoter:                         remoter,
		Repository:                      repo,
		MasternodeID:                    cfg.MasternodePastelID,
		StorageChallengeExpiredDuration: 20 * time.Second,
		PastelClient:                    pastelClient,
	})

	_, err = remoter.RegisterActor(appactor.NewStorageChallengeActor(domainService, store), "storage-challenge")
	if err != nil {
		panic(fmt.Sprintf("coult not register application storage challenge actor: %v", err))
	}

	log.Println("NODE STARTED, INPUT `exit` TO STOP NODE")
	for {
		input, _ := console.ReadLine()
		if input == "exit" {
			return
		}
	}
}
