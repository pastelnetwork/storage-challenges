package actor

import (
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/pastelnetwork/storage-challenges/application/dto"
)

var (
	rootContext *actor.RootContext
	remoter     *remote.Remote
)

func init() {
	system := actor.NewActorSystem()
	rootContext = system.Root
	remoteConfig := remote.Configure("localhost", 8000)
	remoter = remote.NewRemote(system, remoteConfig)
	remoter.Start()
}

func StartStorageChallengeHandler() *actor.PID {
	props := actor.PropsFromProducer(newStorageChallengeActor)
	pid, _ := rootContext.SpawnNamed(props, "storage-challenge")
	return pid
}

func StopActor(pid *actor.PID) {
	rootContext.Stop(pid)
}

func GracefulStop() {
	remoter.Shutdown(true)
}

func newStorageChallengeActor() actor.Actor {
	return &storageChallengeActor{}
}

type storageChallengeActor struct {
	log.Logger //update use another logger provider
	// service    domain.Service
}

func (s *storageChallengeActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *dto.StorageChallengeRequest:
		s.StorageChallenge(context, msg)
	case *dto.VerifyStorageChallengeRequest:
		s.VerifyStorageChallenge(context, msg)
	default:
		log.Printf("Action hot hanled %#v", msg)
		// TODO: response with unhandled notice
	}
}

func (s *storageChallengeActor) StorageChallenge(ctx actor.Context, req *dto.StorageChallengeRequest) (resp *dto.StorageChallengeReply, err error) {
	log.Printf("StorageChallenge handler")
	// validate request body
	es := validateStorageChallengeData(req.GetData(), "Data")
	if err := validationErrorStackWrap(es); err != nil {
		// TODO: send response validation failed to challenger: get address of challenger, response error by ctx.Send(challengeraddressPID, err)
		return &dto.StorageChallengeReply{Data: req.GetData()}, err
	}

	// calling domain service to process bussiness logics
	// service.StorageChallenge()
	// TODO: send response validation failed to challenger and verifyer: get address of challenger and verifyer, send message by ctx.Send(challengeraddressPID, err)
	return &dto.StorageChallengeReply{Data: req.GetData()}, nil
}

func (s *storageChallengeActor) VerifyStorageChallenge(ctx actor.Context, req *dto.VerifyStorageChallengeRequest) (resp *dto.VerifyStorageChallengeReply, err error) {
	log.Printf("VerifyStorageChallenge handler")
	// validate request body
	es := validateStorageChallengeData(req.GetData(), "Data")
	if err := validationErrorStackWrap(es); err != nil {
		return &dto.VerifyStorageChallengeReply{Data: req.GetData()}, err
	}
	// calling domain service to process bussiness logics
	// service.VerifyStorageChallenge()
	// TODO: send response validation failed to challenger and verifyer: get address of challenger and verifyer, send message by ctx.Send(challengeraddressPID, err)
	return &dto.VerifyStorageChallengeReply{Data: req.GetData()}, nil
}
