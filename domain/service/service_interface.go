package service

import (
	"github.com/pastelnetwork/storage-challenges/domain/model"
	appcontext "github.com/pastelnetwork/storage-challenges/utils/context"
)

type StorageChallenge interface {
	GenerateStorageChallenges(ctx appcontext.Context) error
	ProcessStorageChallenge(ctx appcontext.Context, incomingChallengeMessage *model.ChallengeMessages) error
	VerifyStorageChallenge(ctx appcontext.Context, incomingChallengeMessage *model.ChallengeMessages) error
}
