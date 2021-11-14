package repository

import (
	"github.com/pastelnetwork/storage-challenges/domain/model"
	appcontext "github.com/pastelnetwork/storage-challenges/utils/context"
)

type Repository interface {
	GetFilePathFromFileHash(ctx appcontext.Context, fileHash string) (string, error)
	GetSymbolFiles(ctx appcontext.Context) ([]*model.SymbolFile, error)
	GetTopRankedXorDistanceMasternodeToFileHash(ctx appcontext.Context, fileHash string, numberOfChallengeReplicas int, exceptRankedNodeID ...string) ([]*model.XORDistance, error)
	UpsertStorageChallengeMessage(ctx appcontext.Context, storage_challenge_message *model.ChallengeMessage) error

	IncreaseMasternodeTotalChallengesIssued(ctx appcontext.Context, masternodeID string) error
	IncreaseMasternodeTotalChallengesRespondedTo(ctx appcontext.Context, masternodeID string) error
	IncreaseMasternodeTotalChallengesCorrect(ctx appcontext.Context, masternodeID string) error
	IncreaseMasternodeTotalChallengesIncorrect(ctx appcontext.Context, masternodeID string) error
	IncreaseMasternodeTotalChallengesTimeout(ctx appcontext.Context, masternodeID string) error
	IncreasePastelBlockTotalChallengesIssued(ctx appcontext.Context, blockHash string) error
	IncreasePastelBlockTotalChallengesRespondedTo(ctx appcontext.Context, blockHash string) error
	IncreasePastelBlockTotalChallengesCorrect(ctx appcontext.Context, blockHash string) error
	IncreasePastelBlockTotalChallengesIncorrect(ctx appcontext.Context, blockHash string) error
	IncreasePastelBlockTotalChallengesTimeout(ctx appcontext.Context, blockHash string) error
}
