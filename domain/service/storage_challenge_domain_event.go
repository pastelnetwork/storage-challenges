package service

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	actorLog "github.com/AsynkronIT/protoactor-go/log"
	"github.com/pastelnetwork/storage-challenges/application/dto"
	"github.com/pastelnetwork/storage-challenges/domain/model"
)

type verifyStorageChallengeMsg struct {
	VerifierMasterNodesClientPIDs []*actor.PID
	*model.ChallengeMessage
}

func (v *verifyStorageChallengeMsg) String() string {
	return fmt.Sprintf("%#v", v)
}

func (v *verifyStorageChallengeMsg) Reset() {
	v.ChallengeMessage = nil
	v.VerifierMasterNodesClientPIDs = nil
}

func (v *verifyStorageChallengeMsg) ProtoMessage() {}

type processStorageChallengeMsg struct {
	ProcessingMasterNodesClientPID *actor.PID
	*model.ChallengeMessage
}

func (v *processStorageChallengeMsg) String() string {
	return fmt.Sprintf("%#v", v)
}

func (v *processStorageChallengeMsg) Reset() {
	v.ChallengeMessage = nil
	v.ProcessingMasterNodesClientPID = nil
}

func (v *processStorageChallengeMsg) ProtoMessage() {}

type domainActor struct {
}

func (d *domainActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *verifyStorageChallengeMsg:
		d.OnSendVerifyStorageChallengeMessage(context, msg)
	case *processStorageChallengeMsg:
		d.OnSendProcessStorageChallengeMessage(context, msg)
	default:
		log.With(actorLog.String("ACTOR", "DOMAIN ACTOR")).Debug(fmt.Sprintf("Action not handled %#v", msg))
	}
}

func (s *domainActor) OnSendVerifyStorageChallengeMessage(ctx actor.Context, msg *verifyStorageChallengeMsg) {
	for _, verifyingMasternodePID := range msg.VerifierMasterNodesClientPIDs {
		log.Debug(verifyingMasternodePID.String())
		ctx.Send(verifyingMasternodePID, &dto.VerifyStorageChallengeRequest{
			Data: &dto.StorageChallengeData{
				MessageId:                     msg.MessageID,
				MessageType:                   dto.StorageChallengeDataMessageType(dto.StorageChallengeDataMessageType_value[msg.MessageType]),
				ChallengeStatus:               dto.StorageChallengeDataStatus(dto.StorageChallengeDataStatus_value[msg.ChallengeStatus]),
				TimestampChallengeSent:        msg.TimestampChallengeSent,
				TimestampChallengeRespondedTo: msg.TimestampChallengeRespondedTo,
				TimestampChallengeVerified:    0,
				BlockHashWhenChallengeSent:    msg.BlockHashWhenChallengeSent,
				ChallengingMasternodeId:       msg.ChallengingMasternodeID,
				RespondingMasternodeId:        msg.RespondingMasternodeID,
				ChallengeFile: &dto.StorageChallengeDataChallengeFile{
					FileHashToChallenge:      msg.FileHashToChallenge,
					ChallengeSliceStartIndex: int64(msg.ChallengeSliceStartIndex),
					ChallengeSliceEndIndex:   int64(msg.ChallengeSliceEndIndex),
				},
				ChallengeSliceCorrectHash: "",
				ChallengeResponseHash:     msg.ChallengeResponseHash,
				ChallengeId:               msg.ChallengeID,
			},
		})
	}
}

func (s *domainActor) OnSendProcessStorageChallengeMessage(ctx actor.Context, msg *processStorageChallengeMsg) {
	log.Debug(msg.ProcessingMasterNodesClientPID.String())
	ctx.Send(msg.ProcessingMasterNodesClientPID, &dto.StorageChallengeRequest{
		Data: &dto.StorageChallengeData{
			MessageId:                     msg.MessageID,
			MessageType:                   dto.StorageChallengeDataMessageType(dto.StorageChallengeDataMessageType_value[msg.MessageType]),
			ChallengeStatus:               dto.StorageChallengeDataStatus(dto.StorageChallengeDataStatus_value[msg.ChallengeStatus]),
			TimestampChallengeSent:        msg.TimestampChallengeSent,
			TimestampChallengeRespondedTo: 0,
			TimestampChallengeVerified:    0,
			BlockHashWhenChallengeSent:    msg.BlockHashWhenChallengeSent,
			ChallengingMasternodeId:       msg.ChallengingMasternodeID,
			RespondingMasternodeId:        msg.RespondingMasternodeID,
			ChallengeFile: &dto.StorageChallengeDataChallengeFile{
				FileHashToChallenge:      msg.FileHashToChallenge,
				ChallengeSliceStartIndex: int64(msg.ChallengeSliceStartIndex),
				ChallengeSliceEndIndex:   int64(msg.ChallengeSliceEndIndex),
			},
			ChallengeSliceCorrectHash: "",
			ChallengeResponseHash:     "",
			ChallengeId:               msg.ChallengeID,
		},
	})
}
