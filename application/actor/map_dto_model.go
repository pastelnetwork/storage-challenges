package actor

import (
	"github.com/pastelnetwork/storage-challenges/application/dto"
	"github.com/pastelnetwork/storage-challenges/domain/model"
)

func mapChallengeMessage(dto *dto.StorageChallengeData) *model.ChallengeMessages {
	return &model.ChallengeMessages{
		MessageID:                     dto.GetMessageId(),
		MessageType:                   dto.GetMessageType().String(),
		ChallengeStatus:               dto.GetChallengeStatus().String(),
		TimestampChallengeSent:        dto.GetTimestampChallengeSent(),
		TimestampChallengeRespondedTo: dto.GetTimestampChallengeRespondedTo(),
		TimestampChallengeVerified:    dto.GetTimestampChallengeVerified(),
		BlockHashWhenChallengeSent:    dto.GetBlockHashWhenChallengeSent(),
		ChallengingMasternodeId:       dto.GetChallengingMasternodeId(),
		RespondingMasternodeId:        dto.GetRespondingMasternodeId(),
		FileHashToChallenge:           dto.GetChallengeFile().GetFileHashToChallenge(),
		ChallengeSliceStartIndex:      uint64(dto.GetChallengeFile().GetChallengeSliceStartIndex()),
		ChallengeSliceEndIndex:        uint64(dto.GetChallengeFile().GetChallengeSliceEndIndex()),
		ChallengeSliceCorrectHash:     dto.GetChallengeSliceCorrectHash(),
		ChallengeResponseHash:         dto.GetChallengeResponseHash(),
		ChallengeId:                   dto.GetChallengeId(),
	}
}
