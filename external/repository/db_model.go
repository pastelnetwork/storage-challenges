package repository

import (
	"time"

	"github.com/pastelnetwork/storage-challenges/domain/model"
)

type CommonModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ChallengeMessage struct {
	*CommonModel
	*model.ChallengeMessage
}

type Challenge struct {
	*CommonModel
	*model.Challenge
}

type PastelBlock struct {
	*CommonModel
	*model.PastelBlock
}

type Masternode struct {
	*CommonModel
	*model.Masternode
}

type SymbolFile struct {
	*CommonModel
	*model.SymbolFile
}

type XORDistance struct {
	*CommonModel
	*model.XORDistance
}

func (XORDistance) TableName() string {
	return "xor_distances"
}

func mapChallengeMessages(in []*ChallengeMessage) []*model.ChallengeMessage {
	var models = []*model.ChallengeMessage{}
	for _, challengeMessage := range in {
		models = append(models, mapChallengeMessage(challengeMessage))
	}

	return models
}

func mapChallengeMessage(in *ChallengeMessage) *model.ChallengeMessage {
	return in.ChallengeMessage
}

func mapRepoChallengeMessage(in *model.ChallengeMessage) *ChallengeMessage {
	return &ChallengeMessage{CommonModel: &CommonModel{}, ChallengeMessage: in}
}

func mapRepoChallengeMessages(in []*model.ChallengeMessage) []*ChallengeMessage {
	var repos = []*ChallengeMessage{}
	for _, challengeMessage := range in {
		repos = append(repos, mapRepoChallengeMessage(challengeMessage))
	}

	return repos
}

func mapXORDistance(in *XORDistance) *model.XORDistance {
	return in.XORDistance
}

func mapXORDistances(in []*XORDistance) []*model.XORDistance {
	var models = []*model.XORDistance{}
	for _, xorDistance := range in {
		models = append(models, mapXORDistance(xorDistance))
	}

	return models
}

func mapSymbolFile(in *SymbolFile) *model.SymbolFile {
	return in.SymbolFile
}

func mapSymbolFiles(in []*SymbolFile) []*model.SymbolFile {
	var models = []*model.SymbolFile{}
	for _, symbolFile := range in {
		models = append(models, mapSymbolFile(symbolFile))
	}

	return models
}
