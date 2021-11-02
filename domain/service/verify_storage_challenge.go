package service

import (
	"fmt"
	"time"

	actorLog "github.com/AsynkronIT/protoactor-go/log"
	"github.com/pastelnetwork/storage-challenges/domain/model"
	appcontext "github.com/pastelnetwork/storage-challenges/utils/context"
	"github.com/pastelnetwork/storage-challenges/utils/file"
	"github.com/pastelnetwork/storage-challenges/utils/helper"
)

func (s *storageChallenge) VerifyStorageChallenge(ctx appcontext.Context, incomingChallengeMessage *model.ChallengeMessages) error {
	log.With(actorLog.String("ACTOR", "VerifyStorageChallenge")).Debug(incomingChallengeMessage.MessageType)
	if err := s.validateVerifyingStorageChallengeIncommingData(incomingChallengeMessage); err != nil {
		return err
	}

	challengeFileData, err := file.ReadFileIntoMemory(incomingChallengeMessage.FileHashToChallenge)
	if err != nil {
		log.With(actorLog.String("ACTOR", "VerifyStorageChallenge")).Error("could not read file data in to memory", actorLog.String("file.ReadFileIntoMemory", err.Error()))
		return err
	}
	challengeCorrectHash := s.computeHashofFileSlice(challengeFileData, int(incomingChallengeMessage.ChallengeSliceStartIndex), int(incomingChallengeMessage.ChallengeSliceStartIndex))
	messageType := model.MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE
	TimestampChallengeVerified := time.Now().Unix()
	TimeVerifyStorageChallengeInSeconds := helper.ComputeElapsedTimeInSecondsBetweenTwoDatetimes(incomingChallengeMessage.TimestampChallengeSent, TimestampChallengeVerified)
	var challengeStatus string
	if (incomingChallengeMessage.ChallengeResponseHash == challengeCorrectHash) && (TimeVerifyStorageChallengeInSeconds <= float64(s.maxSecondsToRespondToStorageChallenge)) {
		challengeStatus = model.Status_SUCCEEDED
		log.With(actorLog.String("ACTOR", "VerifyStorageChallenge")).Debug("Masternode " + incomingChallengeMessage.RespondingMasternodeId + " correctly responded in " + fmt.Sprint(TimeVerifyStorageChallengeInSeconds) + " seconds to a storage challenge for file " + incomingChallengeMessage.FileHashToChallenge)
	} else if incomingChallengeMessage.ChallengeResponseHash == challengeCorrectHash {
		challengeStatus = model.Status_FAILED_TIMEOUT
		log.With(actorLog.String("ACTOR", "VerifyStorageChallenge")).Debug("Masternode " + incomingChallengeMessage.RespondingMasternodeId + " correctly responded in " + fmt.Sprint(TimeVerifyStorageChallengeInSeconds) + " seconds to a storage challenge for file " + incomingChallengeMessage.FileHashToChallenge + ", but was too slow so failed the challenge anyway!")
	} else {
		challengeStatus = model.Status_FAILED_INCORRECT_RESPONSE
		log.With(actorLog.String("ACTOR", "VerifyStorageChallenge")).Debug("Masternode " + incomingChallengeMessage.RespondingMasternodeId + " failed by incorrectly responding to a storage challenge for file " + incomingChallengeMessage.FileHashToChallenge)
	}

	messageIDInputData := incomingChallengeMessage.ChallengingMasternodeId + incomingChallengeMessage.RespondingMasternodeId + incomingChallengeMessage.FileHashToChallenge + challengeStatus + messageType + incomingChallengeMessage.BlockHashWhenChallengeSent
	messageID := helper.GetHashFromString(messageIDInputData)

	var outgoingChallengeMessage = &model.ChallengeMessages{
		MessageID:                     messageID,
		MessageType:                   model.MessageType_STORAGE_CHALLENGE_VERIFICATION_MESSAGE,
		ChallengeStatus:               challengeStatus,
		TimestampChallengeSent:        incomingChallengeMessage.TimestampChallengeSent,
		TimestampChallengeRespondedTo: incomingChallengeMessage.TimestampChallengeRespondedTo,
		TimestampChallengeVerified:    TimestampChallengeVerified,
		BlockHashWhenChallengeSent:    incomingChallengeMessage.BlockHashWhenChallengeSent,
		ChallengingMasternodeId:       incomingChallengeMessage.ChallengingMasternodeId,
		RespondingMasternodeId:        incomingChallengeMessage.RespondingMasternodeId,
		FileHashToChallenge:           incomingChallengeMessage.FileHashToChallenge,
		ChallengeSliceStartIndex:      incomingChallengeMessage.ChallengeSliceStartIndex,
		ChallengeSliceEndIndex:        incomingChallengeMessage.ChallengeSliceEndIndex,
		ChallengeSliceCorrectHash:     challengeCorrectHash,
		ChallengeResponseHash:         incomingChallengeMessage.ChallengeResponseHash,
		ChallengeId:                   incomingChallengeMessage.ChallengeId,
	}
	s.repository.UpsertStorageChallengeMessage(ctx, outgoingChallengeMessage)
	timeToRespondToStorageChallengeInSeconds := helper.ComputeElapsedTimeInSecondsBetweenTwoDatetimes(incomingChallengeMessage.TimestampChallengeSent, outgoingChallengeMessage.TimestampChallengeRespondedTo)
	log.With(actorLog.String("ACTOR", "VerifyStorageChallenge")).Debug("Masternode " + outgoingChallengeMessage.RespondingMasternodeId + " responded to storage challenge for file hash " + outgoingChallengeMessage.FileHashToChallenge + " in " + fmt.Sprint(timeToRespondToStorageChallengeInSeconds) + " seconds!")

	return nil
}

func (s *storageChallenge) validateVerifyingStorageChallengeIncommingData(incomingChallengeMessage *model.ChallengeMessages) error {
	if incomingChallengeMessage.ChallengeStatus != model.Status_RESPONDED {
		return fmt.Errorf("incorrect status to verifying storage challenge")
	}
	if incomingChallengeMessage.MessageType != model.MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE {
		return fmt.Errorf("incorrect message type to verifying storage challenge")
	}
	return nil
}
