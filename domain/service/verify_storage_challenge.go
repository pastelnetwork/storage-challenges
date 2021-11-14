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

func (s *storageChallenge) VerifyStorageChallenge(ctx appcontext.Context, incomingChallengeMessage *model.ChallengeMessage) error {
	log.With(actorLog.String("ACTOR", "VerifyStorageChallenge")).Debug(incomingChallengeMessage.MessageType)
	if err := s.validateVerifyingStorageChallengeIncommingData(incomingChallengeMessage); err != nil {
		return err
	}

	filePath, err := s.repository.GetFilePathFromFileHash(ctx, incomingChallengeMessage.FileHashToChallenge)
	if err != nil {
		log.With(actorLog.String("ACTOR", "VerifyStorageChallenge")).Error("could not get symbol file path from file hash", actorLog.String("s.repository.GetFilePathFromFileHash", err.Error()))
		return err
	}

	challengeFileData, err := file.ReadFileIntoMemory(filePath)
	if err != nil {
		log.With(actorLog.String("ACTOR", "VerifyStorageChallenge")).Error("could not read file data in to memory", actorLog.String("file.ReadFileIntoMemory", err.Error()))
		return err
	}
	challengeCorrectHash := s.computeHashofFileSlice(challengeFileData, int(incomingChallengeMessage.ChallengeSliceStartIndex), int(incomingChallengeMessage.ChallengeSliceEndIndex))
	messageType := model.MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE
	TimestampChallengeVerified := time.Now().Unix()
	TimeVerifyStorageChallengeInSeconds := helper.ComputeElapsedTimeInSecondsBetweenTwoDatetimes(incomingChallengeMessage.TimestampChallengeSent, TimestampChallengeVerified)
	var challengeStatus string
	var analysisStatus = model.ALALYSIS_STATUS_TIMEOUT
	defer func() {
		s.saveChallengeAnalysis(ctx, incomingChallengeMessage.BlockHashWhenChallengeSent, incomingChallengeMessage.ChallengingMasternodeID, analysisStatus)
	}()

	if (incomingChallengeMessage.ChallengeResponseHash == challengeCorrectHash) && (TimeVerifyStorageChallengeInSeconds <= float64(s.storageChallengeExpiredAsSeconds)) {
		challengeStatus = model.Status_SUCCEEDED
		analysisStatus = model.ANALYSIS_STATUS_CORRECT
		log.With(actorLog.String("ACTOR", "VerifyStorageChallenge")).Debug("Masternode " + incomingChallengeMessage.RespondingMasternodeID + " correctly responded in " + fmt.Sprint(TimeVerifyStorageChallengeInSeconds) + " seconds to a storage challenge for file " + incomingChallengeMessage.FileHashToChallenge)
	} else if incomingChallengeMessage.ChallengeResponseHash == challengeCorrectHash {
		challengeStatus = model.Status_FAILED_TIMEOUT
		analysisStatus = model.ALALYSIS_STATUS_TIMEOUT
		log.With(actorLog.String("ACTOR", "VerifyStorageChallenge")).Debug("Masternode " + incomingChallengeMessage.RespondingMasternodeID + " correctly responded in " + fmt.Sprint(TimeVerifyStorageChallengeInSeconds) + " seconds to a storage challenge for file " + incomingChallengeMessage.FileHashToChallenge + ", but was too slow so failed the challenge anyway!")
	} else {
		challengeStatus = model.Status_FAILED_INCORRECT_RESPONSE
		analysisStatus = model.ALALYSIS_STATUS_INCORRECT
		log.With(actorLog.String("ACTOR", "VerifyStorageChallenge")).Debug("Masternode " + incomingChallengeMessage.RespondingMasternodeID + " failed by incorrectly responding to a storage challenge for file " + incomingChallengeMessage.FileHashToChallenge)
	}

	messageIDInputData := incomingChallengeMessage.ChallengingMasternodeID + incomingChallengeMessage.RespondingMasternodeID + incomingChallengeMessage.FileHashToChallenge + challengeStatus + messageType + incomingChallengeMessage.BlockHashWhenChallengeSent
	messageID := helper.GetHashFromString(messageIDInputData)

	var outgoingChallengeMessage = &model.ChallengeMessage{
		MessageID:                     messageID,
		MessageType:                   model.MessageType_STORAGE_CHALLENGE_VERIFICATION_MESSAGE,
		ChallengeStatus:               challengeStatus,
		TimestampChallengeSent:        incomingChallengeMessage.TimestampChallengeSent,
		TimestampChallengeRespondedTo: incomingChallengeMessage.TimestampChallengeRespondedTo,
		TimestampChallengeVerified:    TimestampChallengeVerified,
		BlockHashWhenChallengeSent:    incomingChallengeMessage.BlockHashWhenChallengeSent,
		ChallengingMasternodeID:       incomingChallengeMessage.ChallengingMasternodeID,
		RespondingMasternodeID:        incomingChallengeMessage.RespondingMasternodeID,
		FileHashToChallenge:           incomingChallengeMessage.FileHashToChallenge,
		ChallengeSliceStartIndex:      incomingChallengeMessage.ChallengeSliceStartIndex,
		ChallengeSliceEndIndex:        incomingChallengeMessage.ChallengeSliceEndIndex,
		ChallengeSliceCorrectHash:     challengeCorrectHash,
		ChallengeResponseHash:         incomingChallengeMessage.ChallengeResponseHash,
		ChallengeID:                   incomingChallengeMessage.ChallengeID,
	}
	if err := s.repository.UpsertStorageChallengeMessage(ctx, outgoingChallengeMessage); err != nil {
		log.With(actorLog.String("ACTOR", "VerifyStorageChallenge")).Error("could not update new storage challenge message in to database", actorLog.String("s.repository.UpsertStorageChallengeMessage", err.Error()))
		return err
	}

	timeToRespondToStorageChallengeInSeconds := helper.ComputeElapsedTimeInSecondsBetweenTwoDatetimes(incomingChallengeMessage.TimestampChallengeSent, outgoingChallengeMessage.TimestampChallengeRespondedTo)
	log.With(actorLog.String("ACTOR", "VerifyStorageChallenge")).Debug("Masternode " + outgoingChallengeMessage.RespondingMasternodeID + " responded to storage challenge for file hash " + outgoingChallengeMessage.FileHashToChallenge + " in " + fmt.Sprint(timeToRespondToStorageChallengeInSeconds) + " seconds!")

	return nil
}

func (s *storageChallenge) validateVerifyingStorageChallengeIncommingData(incomingChallengeMessage *model.ChallengeMessage) error {
	if incomingChallengeMessage.ChallengeStatus != model.Status_RESPONDED {
		return fmt.Errorf("incorrect status to verify storage challenge")
	}
	if incomingChallengeMessage.MessageType != model.MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE {
		return fmt.Errorf("incorrect message type to verify storage challenge")
	}
	return nil
}

func (s *storageChallenge) saveChallengeAnalysis(ctx appcontext.Context, blockHash, challengingMasternodeID string, challengeAnalysisStatus model.ChallengeAnalysisStatus) error {
	switch challengeAnalysisStatus {
	case model.ANALYSYS_STATUS_ISSUED:
		s.repository.IncreaseMasternodeTotalChallengesIssued(ctx, challengingMasternodeID)
		s.repository.IncreasePastelBlockTotalChallengesIssued(ctx, blockHash)
	case model.ANALYSIS_STATUS_RESPONDED_TO:
		s.repository.IncreaseMasternodeTotalChallengesRespondedTo(ctx, challengingMasternodeID)
		s.repository.IncreasePastelBlockTotalChallengesRespondedTo(ctx, blockHash)
	case model.ANALYSIS_STATUS_CORRECT:
		s.repository.IncreaseMasternodeTotalChallengesCorrect(ctx, challengingMasternodeID)
		s.repository.IncreasePastelBlockTotalChallengesCorrect(ctx, blockHash)
	case model.ALALYSIS_STATUS_INCORRECT:
		s.repository.IncreaseMasternodeTotalChallengesIncorrect(ctx, challengingMasternodeID)
		s.repository.IncreasePastelBlockTotalChallengesIncorrect(ctx, blockHash)
	case model.ALALYSIS_STATUS_TIMEOUT:
		s.repository.IncreaseMasternodeTotalChallengesTimeout(ctx, challengingMasternodeID)
		s.repository.IncreasePastelBlockTotalChallengesTimeout(ctx, blockHash)
	}

	return nil
}
