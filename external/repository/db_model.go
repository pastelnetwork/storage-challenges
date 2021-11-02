package repository

import (
	"time"

	"github.com/pastelnetwork/storage-challenges/domain/model"
)

type CommonModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ChallengeMessages struct {
	*CommonModel
	*model.ChallengeMessages
}

type Challenges struct {
	*CommonModel
	*model.Challenges
}

type PastelBlocks struct {
	*CommonModel
	*model.PastelBlocks
}

type Masternodes struct {
	*CommonModel
	*model.Masternodes
}

type SymbolFiles struct {
	*CommonModel
	*model.SymbolFiles
}

type XOR_Distance struct {
	*CommonModel
	*model.XORDistance
}
