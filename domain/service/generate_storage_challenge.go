package service

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	actorLog "github.com/AsynkronIT/protoactor-go/log"
	"github.com/mkmik/argsort"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/storage-challenges/domain/model"
	appcontext "github.com/pastelnetwork/storage-challenges/utils/context"
	"github.com/pastelnetwork/storage-challenges/utils/helper"
)

func (s *storageChallenge) GenerateStorageChallenges(ctx appcontext.Context, challengingMasternodeID string, challengesPerMasternodePerBlock int) error {
	log.With(actorLog.String("ACTOR", "GenerateStorageChallenges")).Debug("Start generating storage challenge")
	curBlockIndex, err := s.pclient.GetBlockCount(ctx)
	if err != nil {
		log.With(actorLog.String("ACTOR", "GenerateStorageChallenges")).Error("could not get block count", actorLog.String("s.pclient.GetBlockCount", err.Error()))
		return err
	}

	curBlockHash, err := s.pclient.GetBlockHash(ctx, curBlockIndex)
	if err != nil {
		log.With(actorLog.String("ACTOR", "GenerateStorageChallenges")).Error("could not get block hash", actorLog.String("s.pclient.getBlockHash", err.Error()))
		return err
	}

	symbolFiles, err := s.repository.GetSymbolFiles(ctx)
	if err != nil {
		log.With(actorLog.String("ACTOR", "GenerateStorageChallenges")).Error("could not get symbol files", actorLog.String("s.repository.GetSymbolFiles", err.Error()))
		return err
	}

	var mapSymbolFileByFileHash = make(map[string]*model.SymbolFile)
	for _, symbolFile := range symbolFiles {
		mapSymbolFileByFileHash[symbolFile.FileHash] = symbolFile
	}

	comparisonStringForFileHashSelection := curBlockHash + challengingMasternodeID
	sliceOfFileHashesToChallenge := getNClosestXORDistanceStringToAGivenComparisonString(challengesPerMasternodePerBlock, comparisonStringForFileHashSelection, _symbolFiles(symbolFiles))

	for idx, symbolFileHash := range sliceOfFileHashesToChallenge {
		challengeDataSize := mapSymbolFileByFileHash[symbolFileHash].FileLengthInBytes

		// selecting n closest node excepting challenger (current node)
		xorDistances, err := s.repository.GetTopRankedXorDistanceMasternodeToFileHash(ctx, symbolFileHash, s.numberOfChallengeReplicas, s.nodeID)
		if err != nil {
			// ignore this file hash
			log.With(actorLog.String("ACTOR", "GenerateStorageChallenges")).Warn(fmt.Sprintf("could not get top %v ranked xor of distance masternodes id to file hash %s", s.numberOfChallengeReplicas, symbolFileHash), actorLog.String("s.repository.GetTopRankedXorDistanceMasternodeToFileHash", err.Error()))
			continue
		}

		comparisonStringForMasternodeSelection := curBlockHash + symbolFileHash + s.nodeID + helper.GetHashFromString(fmt.Sprint(idx))
		respondingMasternodesID := getNClosestXORDistanceStringToAGivenComparisonString(1, comparisonStringForMasternodeSelection, _xorDistances(xorDistances))
		challengeStatus := model.Status_PENDING
		messageType := model.MessageType_STORAGE_CHALLENGE_ISSUANCE_MESSAGE
		challengeSliceStartIndex, challengeSliceEndIndex := getStorageChallengeSliceIndices(uint64(challengeDataSize), symbolFileHash, curBlockHash, challengingMasternodeID)
		messageIDInputData := challengingMasternodeID + respondingMasternodesID[0] + symbolFileHash + challengeStatus + messageType + curBlockHash
		messageID := helper.GetHashFromString(messageIDInputData)
		timestampChallengeSent := time.Now().Unix()
		challengeIDInputData := challengingMasternodeID + respondingMasternodesID[0] + symbolFileHash + fmt.Sprint(challengeSliceStartIndex) + fmt.Sprint(challengeSliceEndIndex) + fmt.Sprint(timestampChallengeSent)
		challengeID := helper.GetHashFromString(challengeIDInputData)
		outgoinChallengMessage := &model.ChallengeMessage{
			MessageID:                     messageID,
			MessageType:                   messageType,
			ChallengeStatus:               challengeStatus,
			TimestampChallengeSent:        time.Now().Unix(),
			TimestampChallengeRespondedTo: 0,
			TimestampChallengeVerified:    0,
			BlockHashWhenChallengeSent:    curBlockHash,
			ChallengingMasternodeID:       challengingMasternodeID,
			RespondingMasternodeID:        respondingMasternodesID[0],
			FileHashToChallenge:           symbolFileHash,
			ChallengeSliceStartIndex:      uint64(challengeSliceStartIndex),
			ChallengeSliceEndIndex:        uint64(challengeSliceEndIndex),
			ChallengeSliceCorrectHash:     "",
			ChallengeResponseHash:         "",
			ChallengeID:                   challengeID,
		}
		err = s.repository.UpsertStorageChallengeMessage(ctx, outgoinChallengMessage)
		if err != nil {
			log.With(actorLog.String("ACTOR", "GenerateStorageChallenges")).Warn(fmt.Sprintf("could not update storage challenge into storage: %v", outgoinChallengMessage), actorLog.String("s.repository.UpsertStorageChallengeMessage", err.Error()))
			continue
		}
		s.sendprocessStorageChallenge(ctx, outgoinChallengMessage)
	}

	return nil
}

func (s *storageChallenge) sendprocessStorageChallenge(ctx appcontext.Context, challengeMessage *model.ChallengeMessage) error {
	masternodes, err := s.pclient.MasterNodesExtra(ctx)
	if err != nil {
		log.With(actorLog.String("ACTOR", "sendprocessStorageChallenge")).Warn("could not get masternode info", actorLog.String("s.pclient.MasterNodesExtra", err.Error()))
		return err
	}

	mapMasternodes := make(map[string]pastel.MasterNode)
	for _, mn := range masternodes {
		mapMasternodes[mn.ExtKey] = mn
	}

	var mn pastel.MasterNode
	var ok bool
	if mn, ok = mapMasternodes[challengeMessage.ChallengingMasternodeID]; !ok {
		err = fmt.Errorf("cannot get masternode info of masternode id %v", challengeMessage.ChallengingMasternodeID)
		log.With(actorLog.String("ACTOR", "sendprocessStorageChallenge")).Warn(fmt.Sprintf("could not get masternode info of %v", challengeMessage.ChallengingMasternodeID), actorLog.String("mapMasternodes[challengeMessage.ChallengingMasternodeID]", err.Error()))
		return err
	}
	processingMasterNodesClientPID := actor.NewPID(mn.ExtAddress, "storage-challenge")

	return s.remoter.Send(ctx, s.domainActorID, &processStorageChallengeMsg{ProcessingMasterNodesClientPID: processingMasterNodesClientPID, ChallengeMessage: challengeMessage})
}

type computingXORDistance interface {
	GetListXORDistanceString() []string
}

type _symbolFiles []*model.SymbolFile

func (s _symbolFiles) GetListXORDistanceString() []string {
	ret := make([]string, len(s))
	for idx, symbolFile := range s {
		ret[idx] = symbolFile.FileHash
	}

	return ret
}

type _xorDistances []*model.XORDistance

func (s _xorDistances) GetListXORDistanceString() []string {
	ret := make([]string, len(s))
	for idx, xorDistance := range s {
		ret[idx] = xorDistance.MasternodeID
	}

	return ret
}

func getNClosestXORDistanceStringToAGivenComparisonString(n int, comparisonString string, sliceOfComputingXORDistance computingXORDistance) []string {
	sliceOfXORDistance := make([]uint64, len(sliceOfComputingXORDistance.GetListXORDistanceString()))
	XORDistanceToComputingStringMap := make(map[uint64]string)
	for idx, currentComputing := range sliceOfComputingXORDistance.GetListXORDistanceString() {
		currentXORDistance := helper.ComputeXorDistanceBetweenTwoStrings(currentComputing, comparisonString)
		sliceOfXORDistance[idx] = currentXORDistance
		XORDistanceToComputingStringMap[currentXORDistance] = currentComputing
	}
	sort.Slice(sliceOfXORDistance, func(i, j int) bool { return sliceOfXORDistance[i] < sliceOfXORDistance[j] })
	sliceOfTopNClosestString := make([]string, n)
	for ii, currentXORDistance := range sliceOfXORDistance {
		if ii < n {
			sliceOfTopNClosestString[ii] = XORDistanceToComputingStringMap[currentXORDistance]
		}
	}
	return sliceOfTopNClosestString
}

func getStorageChallengeSliceIndices(totalDataLengthInBytes uint64, fileHashString string, blockHashString string, challengingMasternodeId string) (int, int) {
	blockHashStringAsInt, _ := strconv.ParseInt(blockHashString, 16, 64)
	blockHashStringAsIntStr := fmt.Sprint(blockHashStringAsInt)
	stepSizeForIndicesStr := blockHashStringAsIntStr[len(blockHashStringAsIntStr)-1:] + blockHashStringAsIntStr[0:1]
	stepSizeForIndices, _ := strconv.ParseUint(stepSizeForIndicesStr, 10, 32)
	stepSizeForIndicesAsInt := int(stepSizeForIndices)
	comparisonString := blockHashString + fileHashString + challengingMasternodeId
	sliceOfXorDistancesOfIndicesToBlockHash := make([]uint64, 0)
	sliceOfIndicesWithStepSize := make([]int, 0)
	totalDataLengthInBytesAsInt := int(totalDataLengthInBytes)
	for j := 0; j <= totalDataLengthInBytesAsInt; j += stepSizeForIndicesAsInt {
		jAsString := fmt.Sprintf("%d", j)
		currentXorDistance := helper.ComputeXorDistanceBetweenTwoStrings(jAsString, comparisonString)
		sliceOfXorDistancesOfIndicesToBlockHash = append(sliceOfXorDistancesOfIndicesToBlockHash, currentXorDistance)
		sliceOfIndicesWithStepSize = append(sliceOfIndicesWithStepSize, j)
	}
	sliceOfSortedIndices := argsort.SortSlice(sliceOfXorDistancesOfIndicesToBlockHash, func(i, j int) bool {
		return sliceOfXorDistancesOfIndicesToBlockHash[i] < sliceOfXorDistancesOfIndicesToBlockHash[j]
	})
	sliceOfSortedIndicesWithStepSize := make([]int, 0)
	for _, currentSortedIndex := range sliceOfSortedIndices {
		sliceOfSortedIndicesWithStepSize = append(sliceOfSortedIndicesWithStepSize, sliceOfIndicesWithStepSize[currentSortedIndex])
	}
	firstTwoSortedIndices := sliceOfSortedIndicesWithStepSize[0:2]
	challengeSliceStartIndex, challengeSliceEndIndex := minMax(firstTwoSortedIndices)
	return challengeSliceStartIndex, challengeSliceEndIndex
}

func minMax(array []int) (int, int) {
	var max int = array[0]
	var min int = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}
