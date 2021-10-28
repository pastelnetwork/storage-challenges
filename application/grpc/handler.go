package grpc

import (
	"context"
	"log"

	"github.com/pastelnetwork/storage-challenges/application/dto"
)

type storageChallenge struct {
	UnimplementedStorageChallengeServer
}

func (s *storageChallenge) StorageChallenge(ctx context.Context, req *dto.StorageChallengeRequest) (*dto.StorageChallengeReply, error) {
	log.Printf("StorageChallenge handler")
	// validate request body
	es := validateStorageChallengeData(req.GetData(), "Data")
	if err := validationErrorStackWrap(es); err != nil {
		return nil, err
	}

	// calling domain service to process bussiness logics
	// domain.service.StorageChallenge()
	return &dto.StorageChallengeReply{Data: req.GetData()}, nil
}

func (s *storageChallenge) VerifyStorageChallenge(ctx context.Context, req *dto.VerifyStorageChallengeRequest) (*dto.VerifyStorageChallengeReply, error) {
	log.Printf("VerifyStorageChallenge handler")
	// validate request body
	es := validateStorageChallengeData(req.GetData(), "Data")
	if err := validationErrorStackWrap(es); err != nil {
		return nil, err
	}
	// calling domain service to process bussiness logics
	// domain.service.VerifyStorageChallenge()
	return &dto.VerifyStorageChallengeReply{}, nil
}

func NewStorageChallengeServer() StorageChallengeServer {
	return &storageChallenge{}
}
