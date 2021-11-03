package actor

import (
	"github.com/pastelnetwork/storage-challenges/application/dto"
	"github.com/pastelnetwork/storage-challenges/domain/model"
)

func mapChallengeMessage(dto *dto.StorageChallengeData) *model.ChallengeMessage {
	return &model.ChallengeMessage{
		MessageID:                     dto.GetMessageId(),
		MessageType:                   dto.GetMessageType().String(),
		ChallengeStatus:               dto.GetChallengeStatus().String(),
		TimestampChallengeSent:        dto.GetTimestampChallengeSent(),
		TimestampChallengeRespondedTo: dto.GetTimestampChallengeRespondedTo(),
		TimestampChallengeVerified:    dto.GetTimestampChallengeVerified(),
		BlockHashWhenChallengeSent:    dto.GetBlockHashWhenChallengeSent(),
		ChallengingMasternodeID:       dto.GetChallengingMasternodeId(),
		RespondingMasternodeID:        dto.GetRespondingMasternodeId(),
		FileHashToChallenge:           dto.GetChallengeFile().GetFileHashToChallenge(),
		ChallengeSliceStartIndex:      uint64(dto.GetChallengeFile().GetChallengeSliceStartIndex()),
		ChallengeSliceEndIndex:        uint64(dto.GetChallengeFile().GetChallengeSliceEndIndex()),
		ChallengeSliceCorrectHash:     dto.GetChallengeSliceCorrectHash(),
		ChallengeResponseHash:         dto.GetChallengeResponseHash(),
		ChallengeID:                   dto.GetChallengeId(),
	}
}
