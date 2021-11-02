package service

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/storage-challenges/domain/model"
	"github.com/pastelnetwork/storage-challenges/external/message"
	"github.com/pastelnetwork/storage-challenges/external/repository"
)

func init() {
}

type storageChallenge struct {
	remoter       *message.Remoter
	repository    repository.Repository
	domainActorID *actor.PID
	nodeID        string
	pclient       pastel.Client
}
type verifyStotageChallengeMsg struct {
	VerifierMasterNodesClientPIDs []*actor.PID
	*model.ChallengeMessages
}

func (v *verifyStotageChallengeMsg) String() string {
	return fmt.Sprintf("%#v", v)
}

func (v *verifyStotageChallengeMsg) Reset() {
	v.ChallengeMessages = nil
	v.VerifierMasterNodesClientPIDs = nil
}

func (v *verifyStotageChallengeMsg) ProtoMessage() {}

type domainActor struct {
}

func (d *domainActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *verifyStotageChallengeMsg:
		d.OnSendVerifyStorageChallengeMessage(context, msg)
	}
}

func NewStorageChallenge(remoter *message.Remoter, repository repository.Repository) StorageChallenge {
	actorID, err := remoter.RegisterActor(&domainActor{}, "domain-service")
	if err != nil {
		panic(err)
	}
	return &storageChallenge{remoter: remoter, repository: repository, domainActorID: actorID}
}
