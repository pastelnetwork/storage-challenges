package service

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	actorLog "github.com/AsynkronIT/protoactor-go/log"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/storage-challenges/domain/model"
	appcontext "github.com/pastelnetwork/storage-challenges/utils/context"
	"github.com/pastelnetwork/storage-challenges/utils/file"
	"github.com/pastelnetwork/storage-challenges/utils/helper"
	"golang.org/x/crypto/sha3"
)

var log *actorLog.Logger

func init() {
	log = actorLog.New(actorLog.DebugLevel, "STORAGE_CHALLENGE")
}

func (s *storageChallenge) ProcessStorageChallenge(ctx appcontext.Context, incomingChallengeMessage *model.ChallengeMessage) error {
	log.With(actorLog.String("ACTOR", "ProcessStorageChallenge")).Debug("Start processing storage challenge")
	if err := s.validateProcessingStorageChallengeIncommingData(incomingChallengeMessage); err != nil {
		return err
	}

	analysisStatus := model.ALALYSIS_STATUS_TIMEOUT

	defer func() {
		s.saveChallengeAnalysis(ctx, incomingChallengeMessage.BlockHashWhenChallengeSent, incomingChallengeMessage.ChallengingMasternodeID, analysisStatus)
	}()

	filePath, err := s.repository.GetFilePathFromFileHash(ctx, incomingChallengeMessage.FileHashToChallenge)
	if err != nil {
		log.With(actorLog.String("ACTOR", "ProcessStorageChallenge")).Error("could not get symbol file path from file hash", actorLog.String("s.repository.GetFilePathFromFileHash", err.Error()))
		return err
	}

	challengeFileData, err := file.ReadFileIntoMemory(filePath)
	if err != nil {
		log.With(actorLog.String("ACTOR", "ProcessStorageChallenge")).Error("could not read file data in to memory", actorLog.String("file.ReadFileIntoMemory", err.Error()))
		return err
	}
	challengeResponseHash := s.computeHashofFileSlice(challengeFileData, int(incomingChallengeMessage.ChallengeSliceStartIndex), int(incomingChallengeMessage.ChallengeSliceEndIndex))
	challengeStatus := model.Status_RESPONDED
	messageType := model.MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE
	messageIDInputData := incomingChallengeMessage.ChallengingMasternodeID + incomingChallengeMessage.RespondingMasternodeID + incomingChallengeMessage.FileHashToChallenge + challengeStatus + messageType + incomingChallengeMessage.BlockHashWhenChallengeSent
	messageID := helper.GetHashFromString(messageIDInputData)
	timestampChallengeRespondedTo := time.Now().Unix()

	var outgoingChallengeMessage = &model.ChallengeMessage{
		MessageID:                     messageID,
		MessageType:                   model.MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE,
		ChallengeStatus:               model.Status_RESPONDED,
		TimestampChallengeSent:        incomingChallengeMessage.TimestampChallengeSent,
		TimestampChallengeRespondedTo: timestampChallengeRespondedTo,
		TimestampChallengeVerified:    0,
		BlockHashWhenChallengeSent:    incomingChallengeMessage.BlockHashWhenChallengeSent,
		ChallengingMasternodeID:       incomingChallengeMessage.ChallengingMasternodeID,
		RespondingMasternodeID:        incomingChallengeMessage.RespondingMasternodeID,
		FileHashToChallenge:           incomingChallengeMessage.FileHashToChallenge,
		ChallengeSliceStartIndex:      incomingChallengeMessage.ChallengeSliceStartIndex,
		ChallengeSliceEndIndex:        incomingChallengeMessage.ChallengeSliceEndIndex,
		ChallengeSliceCorrectHash:     "",
		ChallengeResponseHash:         challengeResponseHash,
		ChallengeID:                   incomingChallengeMessage.ChallengeID,
	}
	if err := s.repository.UpsertStorageChallengeMessage(ctx, outgoingChallengeMessage); err != nil {
		log.With(actorLog.String("ACTOR", "ProcessStorageChallenge")).Error("could not update new storage challenge message in to database", actorLog.String("s.repository.UpsertStorageChallengeMessage", err.Error()))
		return err
	}
	analysisStatus = model.ANALYSIS_STATUS_RESPONDED_TO
	timeToRespondToStorageChallengeInSeconds := helper.ComputeElapsedTimeInSecondsBetweenTwoDatetimes(incomingChallengeMessage.TimestampChallengeSent, outgoingChallengeMessage.TimestampChallengeRespondedTo)
	log.With(actorLog.String("ACTOR", "ProcessStorageChallenge")).Debug("Masternode " + outgoingChallengeMessage.RespondingMasternodeID + " responded to storage challenge for file hash " + outgoingChallengeMessage.FileHashToChallenge + " in " + fmt.Sprint(timeToRespondToStorageChallengeInSeconds) + " seconds!")

	return s.sendVerifyStorageChallenge(ctx, outgoingChallengeMessage)
}

func (s *storageChallenge) validateProcessingStorageChallengeIncommingData(incomingChallengeMessage *model.ChallengeMessage) error {
	if incomingChallengeMessage.ChallengeStatus != model.Status_PENDING {
		return fmt.Errorf("incorrect status to processing storage challenge")
	}
	if incomingChallengeMessage.MessageType != model.MessageType_STORAGE_CHALLENGE_ISSUANCE_MESSAGE {
		return fmt.Errorf("incorrect message type to processing storage challenge")
	}
	return nil
}

func (s *storageChallenge) computeHashofFileSlice(file_data []byte, challenge_slice_start_index int, challenge_slice_end_index int) string {
	challenge_data_slice := file_data[challenge_slice_start_index:challenge_slice_end_index]
	algorithm := sha3.New256()
	algorithm.Write(challenge_data_slice)
	hash_of_data_slice := hex.EncodeToString(algorithm.Sum(nil))
	return hash_of_data_slice
}

func (s *storageChallenge) sendVerifyStorageChallenge(ctx appcontext.Context, challengeMessage *model.ChallengeMessage) error {
	masternodes, err := s.pclient.MasterNodesExtra(ctx)
	if err != nil {
		return err
	}

	mapMasternodes := make(map[string]pastel.MasterNode)
	for _, mn := range masternodes {
		mapMasternodes[mn.ExtKey] = mn
	}

	verifierMasterNodesClientPIDs := []*actor.PID{}
	var mn pastel.MasterNode
	var ok bool
	if mn, ok = mapMasternodes[challengeMessage.RespondingMasternodeID]; !ok {
		return fmt.Errorf("cannot get masternode info of masternode id %v", challengeMessage.RespondingMasternodeID)
	}
	verifierMasterNodesClientPIDs = append(verifierMasterNodesClientPIDs, actor.NewPID(mn.ExtAddress, "storage-challenge"))

	return s.remoter.Send(ctx, s.domainActorID, &verifyStorageChallengeMsg{VerifierMasterNodesClientPIDs: verifierMasterNodesClientPIDs, ChallengeMessage: challengeMessage})
}
