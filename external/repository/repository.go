package repository

import (
	"github.com/pastelnetwork/storage-challenges/domain/model"
	appcontext "github.com/pastelnetwork/storage-challenges/utils/context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repository struct{}

func (r *repository) GetFilePathFromFileHash(ctx appcontext.Context, file_hash_string string) (string, error) {
	db := ctx.GetDBTx()
	var filePath string
	row := db.Table("symbol_files").Where("file_hash = ?", file_hash_string).Select("original_file_path").Row()
	err := row.Scan(&filePath)
	return filePath, err
}

func (r *repository) GetSymbolFiles(ctx appcontext.Context) (list []*model.SymbolFile, err error) {
	db := ctx.GetDBTx()
	var symbolFiles []*SymbolFile
	err = db.Model(&SymbolFile{}).Find(&symbolFiles).Error
	return mapSymbolFiles(symbolFiles), err
}

func (r *repository) GetTopRankedXorDistanceMasternodeToFileHash(ctx appcontext.Context, fileHash string, numberOfChallengeReplicas int, exceptMasternodeIDs ...string) (list []*model.XORDistance, err error) {
	db := ctx.GetDBTx()
	var xorDistances = make([]*XORDistance, 0)
	err = db.Where("symbol_file_hash = ?", fileHash).Not(map[string]interface{}{"masternode_id": exceptMasternodeIDs}).Order("xor_distance ASC").Limit(numberOfChallengeReplicas).Preload("Masternode").Find(&xorDistances).Error
	return mapXORDistances(xorDistances), err
}

func (r *repository) UpsertStorageChallengeMessage(ctx appcontext.Context, challengeMessage *model.ChallengeMessage) error {
	repoChalengeMessage := mapRepoChallengeMessage(challengeMessage)
	db := ctx.GetDBTx()
	return db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "message_id"}}, UpdateAll: true}).Create(repoChalengeMessage).Error
}

func (r *repository) IncreaseMasternodeTotalChallengesIssued(ctx appcontext.Context, masternodeID string) error {
	return ctx.GetDBTx().Model(&Masternode{}).Where("node_id = ?", masternodeID).Update("total_challenges_issued", gorm.Expr("total_challenges_issued + ?", 1)).Error
}

func (r *repository) IncreaseMasternodeTotalChallengesRespondedTo(ctx appcontext.Context, masternodeID string) error {
	return ctx.GetDBTx().Model(&Masternode{}).Where("node_id = ?", masternodeID).Update("total_challenges_responded_to", gorm.Expr("total_challenges_responded_to + ?", 1)).Error
}

func (r *repository) IncreaseMasternodeTotalChallengesCorrect(ctx appcontext.Context, masternodeID string) error {
	return ctx.GetDBTx().Model(&Masternode{}).Where("node_id = ?", masternodeID).Update("total_challenges_correct", gorm.Expr("total_challenges_correct + ?", 1)).Error
}

func (r *repository) IncreaseMasternodeTotalChallengesIncorrect(ctx appcontext.Context, masternodeID string) error {
	return ctx.GetDBTx().Model(&Masternode{}).Where("node_id = ?", masternodeID).Update("total_challenges_incorrect", gorm.Expr("total_challenges_incorrect + ?", 1)).Error
}

func (r *repository) IncreaseMasternodeTotalChallengesTimeout(ctx appcontext.Context, masternodeID string) error {
	return ctx.GetDBTx().Model(&Masternode{}).Where("node_id = ?", masternodeID).Update("total_challenges_timeout", gorm.Expr("total_challenges_timeout + ?", 1)).Error
}

func (r *repository) IncreasePastelBlockTotalChallengesIssued(ctx appcontext.Context, blockHash string) error {
	return ctx.GetDBTx().Model(&PastelBlock{}).Where("block_hash = ?", blockHash).Update("total_challenges_issued", gorm.Expr("total_challenges_issued + ?", 1)).Error
}

func (r *repository) IncreasePastelBlockTotalChallengesRespondedTo(ctx appcontext.Context, blockHash string) error {
	return ctx.GetDBTx().Model(&PastelBlock{}).Where("block_hash = ?", blockHash).Update("total_challenges_responded_to", gorm.Expr("total_challenges_responded_to + ?", 1)).Error
}

func (r *repository) IncreasePastelBlockTotalChallengesCorrect(ctx appcontext.Context, blockHash string) error {
	return ctx.GetDBTx().Model(&PastelBlock{}).Where("block_hash = ?", blockHash).Update("total_challenges_correct", gorm.Expr("total_challenges_correct + ?", 1)).Error
}

func (r *repository) IncreasePastelBlockTotalChallengesIncorrect(ctx appcontext.Context, blockHash string) error {
	return ctx.GetDBTx().Model(&PastelBlock{}).Where("block_hash = ?", blockHash).Update("total_challenges_incorrect", gorm.Expr("total_challenges_incorrect + ?", 1)).Error
}

func (r *repository) IncreasePastelBlockTotalChallengesTimeout(ctx appcontext.Context, blockHash string) error {
	return ctx.GetDBTx().Model(&PastelBlock{}).Where("block_hash = ?", blockHash).Update("total_challenges_timeout", gorm.Expr("total_challenges_timeout + ?", 1)).Error
}

func New() Repository {
	return &repository{}
}
