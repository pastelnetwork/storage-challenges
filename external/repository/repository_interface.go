package repository

import (
	"github.com/pastelnetwork/storage-challenges/domain/model"
	appcontext "github.com/pastelnetwork/storage-challenges/utils/context"
)

type Repository interface {
	GetFilePathFromFileHash(ctx appcontext.Context, file_hash_string string) (string, error) // convert from storage_challenges.GetFilePathFromFileHashUsingDb
	GetXorDistances(ctx appcontext.Context) ([]*model.XOR_Distance, error)                   // convert from storage_challenges.GetXorDistancesFromDb
	GetTopRankedXorDistanceMasternodeToFileHash(ctx appcontext.Context, fileHash string, numberOfChallengeReplicas int, exceptRankedNodeID ...string) ([]*model.XOR_Distance, error)
	FindPendingStorageChallengesByRespondingMasterNodeID(ctx appcontext.Context, responding_masternode_id string) ([]*model.Challenges, error)
	FindRespondedStorageChallengesByChallengingMasterNodeID(ctx appcontext.Context, responding_masternode_id string) ([]*model.Challenges, error)
	FindStorageChallengeInssuanceMessageByChallengeID(ctx appcontext.Context, slice_of_challenge_ids []string) ([]*model.ChallengeMessages, error)
	FindStorageChallengeResponseMessageByChallengeID(ctx appcontext.Context, slice_of_challenge_ids []string) ([]*model.ChallengeMessages, error)
	RemoveMasterNodes(slice_of_pastel_masternode_ids []string) error                                                                                      // convert from storage_challenges.RemoveMasternodesFromDb
	UpsertXorDistances(ctx appcontext.Context, pastel_masternode_ids []string, raptorq_symbol_file_hashes []string, xor_distance_matrix [][]uint64) error // convert from storage_challenges.AddXorDistanceMatrixToDb
	UpsertStorageChallengeMessage(ctx appcontext.Context, storage_challenge_message *model.ChallengeMessages) error                                       // convert from storage_challenges.UpdateDbWithMessage
	UpsertMasternodes(ctx appcontext.Context) error
	UpsertMasternodeStats(ctx appcontext.Context) error    // convert from storage_challenges.UpdateMasternodeStatsInDb
	UpsertBlocks(slice_of_block_hashes []string) error     // convert from storage_challenges.AddBlocksToDb
	UpsertBlockStats(slice_of_block_hashes []string) error // convert from storage_challenges.UpdateBlockStatsInDb
	UpsertFiles(ctx appcontext.Context, slice_of_input_file_paths []string) error
}
