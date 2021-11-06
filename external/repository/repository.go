package repository

import (
	"github.com/pastelnetwork/storage-challenges/domain/model"
	appcontext "github.com/pastelnetwork/storage-challenges/utils/context"
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

func (r *repository) GetXorDistances(ctx appcontext.Context) ([]*model.XORDistance, error) {
	panic("not implemented") // TODO: Implement
}

func (r *repository) GetTopRankedXorDistanceMasternodeToFileHash(ctx appcontext.Context, fileHash string, numberOfChallengeReplicas int, exceptMasternodeIDs ...string) (list []*model.XORDistance, err error) {
	db := ctx.GetDBTx()
	var xorDistances []*XORDistance
	err = db.Preload("Masternode").Where("file_hash = ?", fileHash).Not(map[string]interface{}{"masternode_id": exceptMasternodeIDs}).Order("xor_distance ASC").Limit(numberOfChallengeReplicas).Find(&xorDistances).Error
	return mapXORDistances(xorDistances), err
}

func (r *repository) FindPendingStorageChallengesByRespondingMasterNodeID(ctx appcontext.Context, responding_masternode_id string) ([]*model.Challenge, error) {
	panic("not implemented") // TODO: Implement
}

func (r *repository) FindRespondedStorageChallengesByChallengingMasterNodeID(ctx appcontext.Context, responding_masternode_id string) ([]*model.Challenge, error) {
	panic("not implemented") // TODO: Implement
}

func (r *repository) FindStorageChallengeInssuanceMessageByChallengeID(ctx appcontext.Context, slice_of_challenge_ids []string) ([]*model.ChallengeMessage, error) {
	panic("not implemented") // TODO: Implement
}

func (r *repository) FindStorageChallengeResponseMessageByChallengeID(ctx appcontext.Context, slice_of_challenge_ids []string) ([]*model.ChallengeMessage, error) {
	panic("not implemented") // TODO: Implement
}

func (r *repository) RemoveMasterNodes(slice_of_pastel_masternode_ids []string) error {
	panic("not implemented") // TODO: Implement
}

func (r *repository) UpsertXorDistances(ctx appcontext.Context, pastel_masternode_ids []string, raptorq_symbol_file_hashes []string, xor_distance_matrix [][]uint64) error {
	panic("not implemented") // TODO: Implement
}

func (r *repository) UpsertStorageChallengeMessage(ctx appcontext.Context, challengeMessage *model.ChallengeMessage) error {
	repoChalengeMessage := mapRepoChallengeMessage(challengeMessage)
	db := ctx.GetDBTx()
	return db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "message_id"}}, UpdateAll: true}).Create(repoChalengeMessage).Error
}

func (r *repository) UpsertMasternodes(ctx appcontext.Context) error {
	panic("not implemented") // TODO: Implement
}

func (r *repository) UpsertMasternodeStats(ctx appcontext.Context) error {
	panic("not implemented") // TODO: Implement
}

func (r *repository) UpsertBlocks(slice_of_block_hashes []string) error {
	panic("not implemented") // TODO: Implement
}

func (r *repository) UpsertBlockStats(slice_of_block_hashes []string) error {
	panic("not implemented") // TODO: Implement
}

func (r *repository) UpsertFiles(ctx appcontext.Context, slice_of_input_file_paths []string) error {
	panic("not implemented") // TODO: Implement
}

func New() Repository {
	return &repository{}
}
