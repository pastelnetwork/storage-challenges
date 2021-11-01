package service

import (
	actorLog "github.com/AsynkronIT/protoactor-go/log"
	"github.com/pastelnetwork/storage-challenges/domain/model"
	appcontext "github.com/pastelnetwork/storage-challenges/utils/context"
)

func (s *storageChallenge) VerifyStorageChallenge(ctx appcontext.Context, msg *model.ChallengeMessages) error {
	log.With(actorLog.String("ACTOR", "VerifyStorageChallenge")).Debug(msg.MessageType)
	return nil
}
