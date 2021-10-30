package service

import (
	appcontext "github.com/pastelnetwork/storage-challenges/utils/context"
)

type StorageChallenge interface {
	GenerateStorageChallenges(ctx appcontext.Context) error
	ProcessStorageChallenge(ctx appcontext.Context, request *StorageChallengeData) error
	VerifyStorageChallenge(ctx appcontext.Context, request *StorageChallengeData) error
}
