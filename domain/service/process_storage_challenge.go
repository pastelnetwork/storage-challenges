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

func (s *storageChallenge) ProcessStorageChallenge(ctx appcontext.Context, incommingChallengeMessage *model.ChallengeMessages) error {
	log.With(actorLog.String("ACTOR", "ProcessStorageChallenge")).Debug("Start processing storage challenge")
	challengeFileData, err := file.ReadFileIntoMemory(incommingChallengeMessage.FileHashToChallenge)
	if err != nil {
		log.With(actorLog.String("ACTOR", "ProcessStorageChallenge")).Error("could not read file data in to memory", actorLog.String("file.ReadFileIntoMemory", err.Error()))
		return err
	}
	challengeResponseHash := s.computeHashofFileSlice(challengeFileData, int(incommingChallengeMessage.ChallengeSliceStartIndex), int(incommingChallengeMessage.ChallengeSliceStartIndex))
	challengeStatus := model.Status_RESPONDED
	messageType := model.MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE
	messageIDInputData := incommingChallengeMessage.ChallengingMasternodeId + incommingChallengeMessage.RespondingMasternodeId + incommingChallengeMessage.FileHashToChallenge + challengeStatus + messageType + incommingChallengeMessage.BlockHashWhenChallengeSent
	messageID := helper.GetHashFromString(messageIDInputData)
	timestampChallengeRespondedTo := time.Now().Unix()

	var outGoingChallengeMessage = &model.ChallengeMessages{
		MessageID:                     messageID,
		MessageType:                   model.MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE,
		ChallengeStatus:               model.Status_RESPONDED,
		TimestampChallengeSent:        incommingChallengeMessage.TimestampChallengeSent,
		TimestampChallengeRespondedTo: timestampChallengeRespondedTo,
		TimestampChallengeVerified:    0,
		BlockHashWhenChallengeSent:    incommingChallengeMessage.BlockHashWhenChallengeSent,
		ChallengingMasternodeId:       incommingChallengeMessage.ChallengingMasternodeId,
		RespondingMasternodeId:        incommingChallengeMessage.RespondingMasternodeId,
		FileHashToChallenge:           incommingChallengeMessage.FileHashToChallenge,
		ChallengeSliceStartIndex:      incommingChallengeMessage.ChallengeSliceStartIndex,
		ChallengeSliceEndIndex:        incommingChallengeMessage.ChallengeSliceEndIndex,
		ChallengeSliceCorrectHash:     "",
		ChallengeResponseHash:         challengeResponseHash,
		ChallengeId:                   incommingChallengeMessage.ChallengeId,
	}
	s.repository.UpsertStorageChallengeMessage(ctx, outGoingChallengeMessage)
	timeToRespondToStorageChallengeInSeconds := helper.ComputeElapsedTimeInSecondsBetweenTwoDatetimes(incommingChallengeMessage.TimestampChallengeSent, outGoingChallengeMessage.TimestampChallengeRespondedTo)
	log.With(actorLog.String("ACTOR", "ProcessStorageChallenge")).Debug("Masternode " + outGoingChallengeMessage.RespondingMasternodeId + " responded to storage challenge for file hash " + outGoingChallengeMessage.FileHashToChallenge + " in " + fmt.Sprint(timeToRespondToStorageChallengeInSeconds) + " seconds!")

	return s.sendVerifyStorageChallenge(ctx, outGoingChallengeMessage)
}

func (s *storageChallenge) computeHashofFileSlice(file_data []byte, challenge_slice_start_index int, challenge_slice_end_index int) string {
	challenge_data_slice := file_data[challenge_slice_start_index:challenge_slice_end_index]
	algorithm := sha3.New256()
	algorithm.Write(challenge_data_slice)
	hash_of_data_slice := hex.EncodeToString(algorithm.Sum(nil))
	return hash_of_data_slice
}

func (s *storageChallenge) sendVerifyStorageChallenge(ctx appcontext.Context, challengeMessage *model.ChallengeMessages) error {
	rankedXorDistances, err := s.repository.GetTopRankedXorDistanceMasternodeToFileHash(ctx, challengeMessage.FileHashToChallenge, 6, s.nodeID)
	if err != nil {
		return err
	}

	masternodes, err := s.pclient.MasterNodesExtra(ctx)
	if err != nil {
		return err
	}

	mapMasternodes := make(map[string]pastel.MasterNode)
	for _, mn := range masternodes {
		mapMasternodes[mn.ExtKey] = mn
	}

	verifierMasterNodesClientPIDs := []*actor.PID{}
	for _, xorDistance := range rankedXorDistances {
		var mn pastel.MasterNode
		var ok bool
		if mn, ok = mapMasternodes[xorDistance.MasternodeID]; !ok {
			return fmt.Errorf("cannot get masternode info of masternode id %v", xorDistance.MasternodeID)
		}
		verifierMasterNodesClientPIDs = append(verifierMasterNodesClientPIDs, actor.NewPID(mn.ExtAddress, "domain-service"))
	}

	return s.remoter.Send(ctx, s.domainActorID, &verifyStotageChallengeMsg{VerifierMasterNodesClientPIDs: verifierMasterNodesClientPIDs, ChallengeMessages: challengeMessage})
}
