package service

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pastelnetwork/storage-challenges/application/dto"
)

func (s *domainActor) OnSendVerifyStorageChallengeMessage(ctx actor.Context, msg *verifyStotageChallengeMsg) {
	for _, verifyingMasternodeID := range msg.VerifierMasternodeIDs {
		log.Debug(verifyingMasternodeID)
		// TODO: pastelClient get masternode address from masternodeID, create pid from that address
		ctx.Send(actor.NewPID("localhost:9002", "storage-challenge"), &dto.VerifyStorageChallengeRequest{
			Data: &dto.StorageChallengeData{
				MessageId:                     msg.MessageID,
				MessageType:                   dto.StorageChallengeDataMessageType(dto.StorageChallengeDataStatus_value[msg.MessageType]),
				ChallengeStatus:               dto.StorageChallengeDataStatus(dto.StorageChallengeDataStatus_value[msg.ChallengeStatus]),
				TimestampChallengeSent:        msg.TimestampChallengeSent,
				TimestampChallengeRespondedTo: msg.TimestampChallengeRespondedTo,
				TimestampChallengeVerified:    0,
				BlockHashWhenChallengeSent:    msg.BlockHashWhenChallengeSent,
				ChallengingMasternodeId:       msg.ChallengingMasternodeId,
				RespondingMasternodeId:        msg.RespondingMasternodeId,
				ChallengeFile: &dto.StorageChallengeDataChallengeFile{
					FileHashToChallenge:      msg.FileHashToChallenge,
					ChallengeSliceStartIndex: int64(msg.ChallengeSliceStartIndex),
					ChallengeSliceEndIndex:   int64(msg.ChallengeSliceEndIndex),
				},
				ChallengeSliceCorrectHash: "",
				ChallengeResponseHash:     msg.ChallengeResponseHash,
				ChallengeId:               msg.ChallengeId,
			},
		})
	}
}
