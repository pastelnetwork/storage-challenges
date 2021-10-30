package service

import (
	actorLog "github.com/AsynkronIT/protoactor-go/log"
	appcontext "github.com/pastelnetwork/storage-challenges/utils/context"
)

func (s *storageChallenge) VerifyStorageChallenge(ctx appcontext.Context, request *StorageChallengeData) error {
	log.With(actorLog.String("ACTOR", "VerifyStorageChallenge")).Debug(request.MessageType.String())
	return nil
}
