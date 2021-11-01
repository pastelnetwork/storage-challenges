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
	var dbTx = s.db.GetDB().Begin()
	var commit bool
	defer func() {
		if !commit {
			dbTx.Rollback()
			return
		}
		dbTx.Commit()
	}()

	var ctx = appcontext.FromContext(context.Background()).WithActorContext(actorCtx).WithDBTx(dbTx)
	switch msg := actorCtx.Message().(type) {
	case *dto.StorageChallengeRequest:
		_, err := s.StorageChallenge(ctx, msg)
		if err == nil {
			commit = true
		}
	case *dto.VerifyStorageChallengeRequest:
		_, err := s.VerifyStorageChallenge(ctx, msg)
		if err == nil {
			commit = true
		}
	default:
		log.Printf("Action not hanled %#v", msg)
		// TODO: response with unhandled notice
	}
}

func (s *storageChallengeActor) StorageChallenge(ctx appcontext.Context, req *dto.StorageChallengeRequest) (resp *dto.StorageChallengeReply, err error) {
	log.Printf("StorageChallenge handler")
	// validate request body
	es := validateStorageChallengeData(req.GetData(), "Data")
	if err := validationErrorStackWrap(es); err != nil {
		// TODO: send response validation failed to challenger: get address of challenger, response error by ctx.Send(challengeraddressPID, err)
		return &dto.StorageChallengeReply{Data: req.GetData()}, err
	}

	// calling domain service to process bussiness logics
	err = s.service.ProcessStorageChallenge(ctx, mapChallengeMessage(req.GetData()))
	// TODO: send response validation failed to challenger and verifyer: get address of challenger and verifyer, send message by ctx.Send(challengeraddressPID, err)
	return &dto.StorageChallengeReply{Data: req.GetData()}, err
}

func (s *storageChallengeActor) VerifyStorageChallenge(ctx appcontext.Context, req *dto.VerifyStorageChallengeRequest) (resp *dto.VerifyStorageChallengeReply, err error) {
	log.Printf("VerifyStorageChallenge handler")
	// validate request body
	es := validateStorageChallengeData(req.GetData(), "Data")
	if err := validationErrorStackWrap(es); err != nil {
		return &dto.VerifyStorageChallengeReply{Data: req.GetData()}, err
	}
	// calling domain service to process bussiness logics
	err = s.service.VerifyStorageChallenge(ctx, mapChallengeMessage(req.GetData()))
	// TODO: send response validation failed to challenger and verifyer: get address of challenger and verifyer, send message by ctx.Send(challengeraddressPID, err)
	return &dto.VerifyStorageChallengeReply{Data: req.GetData()}, err
}
