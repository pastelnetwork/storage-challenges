package service

import (
	"github.com/pastelnetwork/storage-challenges/external/message"
	"github.com/pastelnetwork/storage-challenges/external/repository"
)

func init() {
}

type storageChallenge struct {
	remoter    *message.Remoter
	repository repository.Repository
}

func NewStorageChallenge(remoter *message.Remoter, repository repository.Repository) StorageChallenge {
	return &storageChallenge{remoter: remoter, repository: repository}
}
