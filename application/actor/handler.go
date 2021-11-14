package actor

import (
	"context"
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pastelnetwork/storage-challenges/application/dto"
	"github.com/pastelnetwork/storage-challenges/domain/service"
	"github.com/pastelnetwork/storage-challenges/external/storage"
	appcontext "github.com/pastelnetwork/storage-challenges/utils/context"
)

func NewStorageChallengeActor(domainService service.StorageChallenge, db storage.Store) actor.Actor {
	return &storageChallengeActor{service: domainService, db: db}
}

type storageChallengeActor struct {
	service service.StorageChallenge
	db      storage.Store
}

func (s *storageChallengeActor) Receive(actorCtx actor.Context) {
	// Begin transaction, inject to context to go through main process

	var ctx = appcontext.FromContext(context.Background()).WithActorContext(actorCtx).WithDBTx(s.db.GetDB())
	switch msg := actorCtx.Message().(type) {
	case *dto.GenerateStorageChallengeRequest:
		s.GenerateStorageChallenges(ctx, msg)
	case *dto.StorageChallengeRequest:
		s.StorageChallenge(ctx, msg)
	case *dto.VerifyStorageChallengeRequest:
		s.VerifyStorageChallenge(ctx, msg)
	default:
		log.Printf("Action not handled %#v", msg)
		// TODO: response with unhandled notice
	}
}

func (s *storageChallengeActor) GenerateStorageChallenges(ctx appcontext.Context, req *dto.GenerateStorageChallengeRequest) (resp *dto.GenerateStorageChallengeReply, err error) {
	log.Println("GenerateStorageChallenge handler")
	// validate request body
	es := validateGenerateStorageChallengeData(req)
	if err := validationErrorStackWrap(es); err != nil {
		log.Printf("[GenerateStorageChallenge][Validation Error] %s", es)
		return &dto.GenerateStorageChallengeReply{}, err
	}

	// calling domain service to process bussiness logics
	err = s.service.GenerateStorageChallenges(ctx, req.GetCurrentBlockHash(), req.GetChallengingMasternodeId(), int(req.GetChallengesPerMasternodePerBlock()))
	return &dto.GenerateStorageChallengeReply{}, err
}

func (s *storageChallengeActor) StorageChallenge(ctx appcontext.Context, req *dto.StorageChallengeRequest) (resp *dto.StorageChallengeReply, err error) {
	log.Println("StorageChallenge handler")
	// validate request body
	es := validateStorageChallengeData(req.GetData(), "Data")
	if err := validationErrorStackWrap(es); err != nil {
		log.Printf("[ProcessStorageChallenge][Validation Error] %s", es)
		return &dto.StorageChallengeReply{Data: req.GetData()}, err
	}

	// calling domain service to process bussiness logics
	err = s.service.ProcessStorageChallenge(ctx, mapChallengeMessage(req.GetData()))
	return &dto.StorageChallengeReply{Data: req.GetData()}, err
}

func (s *storageChallengeActor) VerifyStorageChallenge(ctx appcontext.Context, req *dto.VerifyStorageChallengeRequest) (resp *dto.VerifyStorageChallengeReply, err error) {
	log.Println("VerifyStorageChallenge handler")
	// validate request body
	es := validateStorageChallengeData(req.GetData(), "Data")
	if err := validationErrorStackWrap(es); err != nil {
		log.Printf("[VerifyStorageChallenge][Validation Error] %s", es)
		return &dto.VerifyStorageChallengeReply{Data: req.GetData()}, err
	}
	// calling domain service to process bussiness logics
	err = s.service.VerifyStorageChallenge(ctx, mapChallengeMessage(req.GetData()))
	return &dto.VerifyStorageChallengeReply{Data: req.GetData()}, err
}
