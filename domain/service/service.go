package service

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/storage-challenges/external/message"
	"github.com/pastelnetwork/storage-challenges/external/repository"
)

type Config struct {
	Remoter                               *message.Remoter
	Repository                            repository.Repository
	MasterNodeID                          string
	PastelClient                          pastel.Client
	MaxSecondsToRespondToStorageChallenge int64
}

type storageChallenge struct {
	remoter                               *message.Remoter
	repository                            repository.Repository
	domainActorID                         *actor.PID
	nodeID                                string
	pclient                               pastel.Client
	maxSecondsToRespondToStorageChallenge int64
}

func NewStorageChallenge(cfg Config) StorageChallenge {
	actorID, err := cfg.Remoter.RegisterActor(&domainActor{}, "domain-service")
	if err != nil {
		panic(err)
	}
	return &storageChallenge{remoter: cfg.Remoter, repository: cfg.Repository, domainActorID: actorID, nodeID: cfg.MasterNodeID, pclient: cfg.PastelClient}
}
