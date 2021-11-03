package service

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pastelnetwork/storage-challenges/application/dto"
	"github.com/pastelnetwork/storage-challenges/domain/model"
)

type verifyStotageChallengeMsg struct {
	VerifierMasterNodesClientPIDs []*actor.PID
	*model.ChallengeMessage
}

func (v *verifyStotageChallengeMsg) String() string {
	return fmt.Sprintf("%#v", v)
}

func (v *verifyStotageChallengeMsg) Reset() {
	v.ChallengeMessage = nil
	v.VerifierMasterNodesClientPIDs = nil
}

func (v *verifyStotageChallengeMsg) ProtoMessage() {}

type domainActor struct {
}

func (d *domainActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *verifyStotageChallengeMsg:
		d.OnSendVerifyStorageChallengeMessage(context, msg)
	}
}

func (s *domainActor) OnSendVerifyStorageChallengeMessage(ctx actor.Context, msg *verifyStotageChallengeMsg) {
	for _, verifyingMasternodePID := range msg.VerifierMasterNodesClientPIDs {
		log.Debug(verifyingMasternodePID.String())
		ctx.Send(verifyingMasternodePID, &dto.VerifyStorageChallengeRequest{
			Data: &dto.StorageChallengeData{
				MessageId:                     msg.MessageID,
				MessageType:                   dto.StorageChallengeDataMessageType(dto.StorageChallengeDataStatus_value[msg.MessageType]),
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
