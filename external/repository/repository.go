package repository

import (
	"fmt"
	"strings"

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

func (r *repository) GetXorDistances(ctx appcontext.Context) ([]*model.XORDistance, error) {
	panic("not implemented") // TODO: Implement
}

func (r *repository) GetTopRankedXorDistanceMasternodeToFileHash(ctx appcontext.Context, numberOfChallengeReplicas int, exceptMasternodeIDs ...string) ([]*model.XORDistance, error) {
	var queryStatement, queryStatementPrepared string
	if len(exceptMasternodeIDs) == 0 {
		queryStatement = "SELECT xor_distance_id, masternode_id, file_hash, xor_distance FROM (SELECT *, RANK() OVER (PARTITION BY file_hash ORDER BY xor_distance ASC) as rnk FROM (SELECT * FROM xor_distances)) WHERE rnk <= %v"
		queryStatementPrepared = fmt.Sprintf(queryStatement, numberOfChallengeReplicas)
	} else {
		queryStatement = "SELECT xor_distance_id, masternode_id, file_hash, xor_distance FROM (SELECT *, RANK() OVER (PARTITION BY file_hash ORDER BY xor_distance ASC) as rnk FROM (SELECT * FROM xor_distances WHERE masternode_id NOT IN (%s))) WHERE rnk <= %v"
		queryStatementPrepared = fmt.Sprintf(queryStatement, strings.Join(exceptMasternodeIDs, ","), numberOfChallengeReplicas)
	}
	db := ctx.GetDBTx()

	var xorDistances []*XORDistance
	err := db.Model(&model.XORDistance{}).Raw(queryStatementPrepared).Scan(&xorDistances).Error
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

func (repository) UpsertFiles(ctx appcontext.Context, slice_of_input_file_paths []string) error {
	panic("not implemented") // TODO: Implement
}
