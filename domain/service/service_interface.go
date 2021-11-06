package service

import (
	"github.com/pastelnetwork/storage-challenges/domain/model"
	appcontext "github.com/pastelnetwork/storage-challenges/utils/context"
)

type StorageChallenge interface {
	GenerateStorageChallenges(ctx appcontext.Context, challengingMasternodeID string, challengesPerNodePerBlock int) error
	ProcessStorageChallenge(ctx appcontext.Context, incomingChallengeMessage *model.ChallengeMessage) error
	VerifyStorageChallenge(ctx appcontext.Context, incomingChallengeMessage *model.ChallengeMessage) error
}
