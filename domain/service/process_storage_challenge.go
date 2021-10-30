package service

import (
	actorLog "github.com/AsynkronIT/protoactor-go/log"
	"github.com/pastelnetwork/storage-challenges/application/dto"
	"github.com/pastelnetwork/storage-challenges/external/message"
	appcontext "github.com/pastelnetwork/storage-challenges/utils/context"
)

var log *actorLog.Logger

func init() {
	log = actorLog.New(actorLog.DebugLevel, "STORAGE_CHALLENGE")
}

func (s *storageChallenge) ProcessStorageChallenge(ctx appcontext.Context, request *StorageChallengeData) error {
	log.With(actorLog.String("ACTOR", "ProcessStorageChallenge")).Debug(request.MessageType.String())
	request.MessageType = dto.MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE

	return s.sendVerifyStorageChallenge(ctx, request)
}

func (s *storageChallenge) sendVerifyStorageChallenge(ctx appcontext.Context, data *StorageChallengeData) error {
	return s.remoter.SendMany(ctx, []message.ActorProperties{
		{
			Address: "localhost:9002",
			Name:    "storage-challenge",
			Kind:    "storage-challenge",
		},
		{
			Address: "localhost:9003",
			Name:    "storage-challenge",
			Kind:    "storage-challenge",
		},
	}, &dto.VerifyStorageChallengeRequest{Data: data})
}
