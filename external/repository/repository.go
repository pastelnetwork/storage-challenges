package repository

import (
	"github.com/pastelnetwork/storage-challenges/domain/model"
	appcontext "github.com/pastelnetwork/storage-challenges/utils/context"
)

type repository struct{}

func (r *repository) GetFilePathFromFileHash(ctx appcontext.Context, file_hash_string string) (string, error) {
	db := ctx.GetDBTx()
	var file_path string
	row := db.Table("symbol_files").Where("file_hash = ?", file_hash_string).Select("original_file_path").Row()
	return file_path, row.Scan(&file_path)
}

func (r *repository) GetXorDistances(ctx appcontext.Context) ([]*model.XOR_Distance, error) {
	panic("not implemented") // TODO: Implement
}

func (r *repository) GetTopRankedXorDistanceMasternodeToFileHash(ctx appcontext.Context, number_of_challenge_replicas int) ([]*model.XOR_Distance, error) {
	panic("not implemented") // TODO: Implement
}

func (r *repository) FindPendingStorageChallengesByRespondingMasterNodeID(ctx appcontext.Context, responding_masternode_id string) ([]*model.Challenges, error) {
	panic("not implemented") // TODO: Implement
}

func (r *repository) FindRespondedStorageChallengesByChallengingMasterNodeID(ctx appcontext.Context, responding_masternode_id string) ([]*model.Challenges, error) {
	panic("not implemented") // TODO: Implement
}

func (r *repository) FindStorageChallengeInssuanceMessageByChallengeID(ctx appcontext.Context, slice_of_challenge_ids []string) ([]*model.ChallengeMessages, error) {
	panic("not implemented") // TODO: Implement
}

func (r *repository) FindStorageChallengeResponseMessageByChallengeID(ctx appcontext.Context, slice_of_challenge_ids []string) ([]*model.ChallengeMessages, error) {
	panic("not implemented") // TODO: Implement
}

func (r *repository) RemoveMasterNodes(slice_of_pastel_masternode_ids []string) error {
	panic("not implemented") // TODO: Implement
}

func (r *repository) UpsertXorDistances(ctx appcontext.Context, pastel_masternode_ids []string, raptorq_symbol_file_hashes []string, xor_distance_matrix [][]uint64) error {
	panic("not implemented") // TODO: Implement
}

func (r *repository) UpsertStorageChallengeMessage(ctx appcontext.Context, storage_challenge_message *model.ChallengeMessages) error {
	panic("not implemented") // TODO: Implement
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
