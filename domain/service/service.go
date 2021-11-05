package service

import (
	"fmt"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/storage-challenges/external/message"
	"github.com/pastelnetwork/storage-challenges/external/repository"
)

type Config struct {
	Remoter                         *message.Remoter
	Repository                      repository.Repository
	MasternodeID                    string
	PastelClient                    pastel.Client
	StorageChallengeExpiredDuration time.Duration
	NumberOfChallengeReplicas       int
}

type storageChallenge struct {
	remoter                          *message.Remoter
	repository                       repository.Repository
	domainActorID                    *actor.PID
	nodeID                           string
	pclient                          pastel.Client
	storageChallengeExpiredAsSeconds int64
	numberOfChallengeReplicas        int
}

func NewStorageChallenge(cfg Config) StorageChallenge {
	actorID, err := cfg.Remoter.RegisterActor(&domainActor{}, "domain-service")
	if err != nil {
		panic(fmt.Sprintf("could not register domain actor: %v", err))
	}
	return &storageChallenge{remoter: cfg.Remoter, repository: cfg.Repository, domainActorID: actorID, nodeID: cfg.MasternodeID, pclient: cfg.PastelClient, storageChallengeExpiredAsSeconds: int64(cfg.StorageChallengeExpiredDuration.Seconds())}
}
