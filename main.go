package main

import (
	"log"
	"net"

	appgrpc "github.com/pastelnetwork/storage-challenges/application/grpc"
	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer(grpc.EmptyServerOption{})
	appgrpc.RegisterStorageChallengeServer(server, appgrpc.NewStorageChallengeServer())
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("server started")
	server.Serve(l)
	defer server.GracefulStop()
}

// import (
// 	"fmt"
// 	"log"
// 	"math"
// 	"path/filepath"
// 	"time"

// 	"github.com/pastelnetwork/storage-challenges/config"
// 	"github.com/pastelnetwork/storage-challenges/storagechallenges"
// 	"github.com/pastelnetwork/storage-challenges/utils/file"
// 	"github.com/pastelnetwork/storage-challenges/utils/helper"
// 	"gorm.io/driver/sqlite"
// 	"gorm.io/gorm"
// )

// func main() {

// 	db, err := gorm.Open(sqlite.Open("go_pastel_storage_challenges.sqlite"), &gorm.Config{})
// 	if err != nil {
// 		panic("failed to connect database")
// 	}

// 	type Masternodes struct {
// 		gorm.Model
// 		Masternode_id                         string `gorm:"primaryKey;unique"`
// 		Masternode_ip_address                 string
// 		Total_challenges_issued               uint
// 		Total_challenges_responded_to         uint
// 		Total_challenges_correct              uint
// 		Total_challenges_incorrect            uint
// 		Total_challenges_correct_but_too_slow uint
// 		Total_challenges_never_responded_to   uint
// 		Challenge_response_success_rate_pct   float32
// 	}

// 	type PastelBlocks struct {
// 		gorm.Model
// 		Block_hash                            string `gorm:"primaryKey;unique"`
// 		Block_number                          uint
// 		Total_challenges_issued               uint
// 		Total_challenges_responded_to         uint
// 		Total_challenges_correct              uint
// 		Total_challenges_incorrect            uint
// 		Total_challenges_correct_but_too_slow uint
// 		Total_challenges_never_responded_to   uint
// 		Challenge_response_success_rate_pct   float32
// 	}

// 	type SymbolFiles struct {
// 		gorm.Model
// 		File_hash                 string `gorm:"primaryKey;unique"`
// 		File_length_in_bytes      uint
// 		Total_challenges_for_file uint
// 		Original_file_path        string
// 	}

// 	type XOR_Distance struct {
// 		gorm.Model
// 		Xor_distance_id string `gorm:"primaryKey;unique"`
// 		Masternode_id   string
// 		File_hash       string
// 		Xor_distance    uint64
// 		SymbolFiles     SymbolFiles `gorm:"foreignKey:Masternode_id"`
// 		Masternodes     Masternodes `gorm:"foreignKey:File_hash"`
// 	}

// 	type Challenges struct {
// 		gorm.Model
// 		Challenge_id                       string `gorm:"primaryKey;unique"`
// 		Challenge_status                   string
// 		Datetime_challenge_sent            string
// 		Datetime_challenge_responded_to    string
// 		Datetime_challenge_verified        string
// 		Block_hash_when_challenge_sent     string
// 		Challenge_response_time_in_seconds float64
// 		Challenging_masternode_id          string
// 		Responding_masternode_id           string
// 		File_hash_to_challenge             string
// 		Challenge_slice_start_index        uint64
// 		Challenge_slice_end_index          uint64
// 		Challenge_slice_correct_hash       string
// 		Challenge_response_hash            string
// 		PastelBlocks                       PastelBlocks `gorm:"foreignKey:Block_hash_when_challenge_sent"`
// 		SymbolFiles                        SymbolFiles  `gorm:"foreignKey:File_hash_to_challenge"`
// 		Masternodes                        Masternodes  `gorm:"foreignKey:Challenging_masternode_id; foreignKey:Responding_masternode_id"`
// 	}

// 	type ChallengeMessages struct {
// 		gorm.Model
// 		Message_id                      string `gorm:"primaryKey;unique"`
// 		Message_type                    string
// 		Challenge_status                string
// 		Datetime_challenge_sent         string
// 		Datetime_challenge_responded_to string
// 		Datetime_challenge_verified     string
// 		Block_hash_when_challenge_sent  string
// 		Challenging_masternode_id       string
// 		Responding_masternode_id        string
// 		File_hash_to_challenge          string
// 		Challenge_slice_start_index     uint64
// 		Challenge_slice_end_index       uint64
// 		Challenge_slice_correct_hash    string
// 		Challenge_response_hash         string
// 		Challenge_id                    string
// 		PastelBlocks                    PastelBlocks `gorm:"foreignKey:Block_hash_when_challenge_sent"`
// 		SymbolFiles                     SymbolFiles  `gorm:"foreignKey:File_hash_to_challenge"`
// 		Masternodes                     Masternodes  `gorm:"foreignKey:Challenging_masternode_id; foreignKey:Responding_masternode_id"`
// 		Challenges                      Challenges   `gorm:"foreignKey:Challenge_id"`
// 	}

// 	db.AutoMigrate(
// 		&Masternodes{},
// 		&PastelBlocks{},
// 		&SymbolFiles{},
// 		&XOR_Distance{},
// 		&Challenges{},
// 		&ChallengeMessages{},
// 	)

// 	run_tests := false

// 	if run_tests {
// 		sample_filepath := "D:\\dupe_detection_downloaded_images\\true_original_images\\2ab79ec61e.jpg"
// 		fmt.Println((sample_filepath))

// 		file_hash, _ := file.GetHashFromFilePath(sample_filepath)
// 		fmt.Println((file_hash))

// 		sample_string := "hello"
// 		hash_of_string := helper.GetHashFromString(sample_string)
// 		fmt.Println(hash_of_string)

// 		sample_string1 := "abc"
// 		sample_string2 := "cde"
// 		sample_xor_distance := helper.ComputeXorDistanceBetweenTwoStrings(sample_string1, sample_string2)
// 		fmt.Println(sample_xor_distance)

// 		sample_fake_masternode_id := helper.GenerateFakePastelMnID()
// 		fmt.Println(sample_fake_masternode_id)

// 		sample_number_of_fake_block_hashes := 20
// 		sample_slice_of_fake_block_hashes := helper.GenerateFakeBlockHashes(sample_number_of_fake_block_hashes)
// 		fmt.Println(sample_slice_of_fake_block_hashes)

// 		sample_start_datetime_string := "2021-09-29T11:45:26.371Z"
// 		sample_end_datetime_string := time.Now().Format(time.RFC3339)
// 		sample_duration_in_seconds := helper.ComputeElapsedTimeInSecondsBetweenTwoDatetimes(sample_start_datetime_string, sample_end_datetime_string)
// 		fmt.Println(fmt.Sprint(sample_duration_in_seconds))

// 		const path_to_raptorq_files = "D:\\je_golang_storage_challenge_code\\sample_raptorq_symbol_files\\"
// 		sample_slice_of_raptorq_file_hashes := storagechallenges.GetRaptorqFileHashes(path_to_raptorq_files)
// 		fmt.Println(sample_slice_of_raptorq_file_hashes)

// 		sample_slice_of_masternode_ids := []string{"jXYqNKXL9otaCCojGNoeA7zv1UZYTigAnUwDQyonh8UqVEjKJcuagRGfMbanaWyVKCIKKFy7FSIaqVzZFCudQE", "jXamyHN3fe8Y8aLMdnVnMoGHJ8Y6nxVrLbwPpXwiarYxw5BROiKkZjSeVBqS6NrwluLrJXnVBhmg1hZLRsvtLq", "jXWs6vA0vewPohb1u0IdX9qVisE56J0XRfOG3ItVp98bq22p6r90QtTlitT9w5FLvC5fvpspJY3utgfxAGBKbp", "jXMZ1jP38I88C7YraMHzEw4lhGxHWnw5gkoZitRBfGG5hDtMpopdGMfof0FI2ruEsXIwAgcFcu1FUeOX30yHvf", "jXFGrrpRuxUczmblectl9Lk3C7fWJrXCwoYvp1kBxk1SvLlbD3nh0ZijLBE0Ut4t6RIGcWirLnxv32DhcYLDaZ"}
// 		//sample_number_of_fake_masternode_ids := 5
// 		//sample_slice_of_masternode_ids := storagechallenges.GetMasternodeIds(sample_number_of_fake_masternode_ids)
// 		fmt.Println(sample_slice_of_masternode_ids)

// 		sample_xor_distance_matrix := storagechallenges.ComputeMasternodeIdToFileHashXorDistanceMatrix(sample_slice_of_masternode_ids, sample_slice_of_raptorq_file_hashes)

// 		storagechallenges.AddXorDistanceMatrixToDb(sample_slice_of_masternode_ids, sample_slice_of_raptorq_file_hashes, sample_xor_distance_matrix)
// 		fmt.Println("Added XOR distances to database!")

// 		//storagechallenges.GetXorDistancesFromDb()
// 		number_of_storage_replicas := 5
// 		sample_masternode_to_file_hash_responsibility_slice_of_structs := storagechallenges.DetermineWhichMasternodesAreResponsibleForWhichFileHashes(number_of_storage_replicas)
// 		fmt.Println(sample_masternode_to_file_hash_responsibility_slice_of_structs)

// 		sample_path_to_raptorq_files := "D:\\je_golang_storage_challenge_code\\sample_raptorq_symbol_files\\"
// 		sample_slice_of_file_paths := storagechallenges.GetFilePathsFromFolder(sample_path_to_raptorq_files)
// 		storagechallenges.AddFilesToDb(sample_slice_of_file_paths)

// 		samp_slice_of_masternode_ids, samp_slice_of_file_hashes := storagechallenges.GetCurrentListsOfMasternodeIdsAndFileHashesFromDb()
// 		sample_slice_of_file_paths2 := storagechallenges.GetSliceOfFilePathsFromSliceOfFileHashes(samp_slice_of_file_hashes)
// 		fmt.Println(sample_slice_of_file_paths2)

// 		fmt.Println(samp_slice_of_masternode_ids)
// 		fmt.Println(samp_slice_of_file_hashes)

// 		sample_file_hash := samp_slice_of_file_hashes[1]
// 		sample_n := 3
// 		sample_top_n_closest_masternode_ids := storagechallenges.GetNClosestMasternodeIdsToAGivenFileHashUsingDb(sample_n, sample_file_hash)
// 		fmt.Println(sample_top_n_closest_masternode_ids)

// 		sample_comparison_string := "test123"
// 		sample_top_n_closest_masternode_ids2 := storagechallenges.GetNClosestMasternodeIdsToAGivenComparisonString(sample_n, sample_comparison_string, samp_slice_of_masternode_ids)
// 		fmt.Println(sample_top_n_closest_masternode_ids2)

// 		sample_top_n_closest_file_hashes := storagechallenges.GetNClosestFileHashesToAGivenComparisonString(sample_n, sample_comparison_string, samp_slice_of_file_hashes)
// 		fmt.Println(sample_top_n_closest_file_hashes)

// 		sample_total_data_length_in_bytes := uint64(50000)
// 		sample_file_hash_string := samp_slice_of_file_hashes[0]
// 		sample_block_hash_string := sample_slice_of_fake_block_hashes[0]
// 		sample_challenging_masternode_id := samp_slice_of_masternode_ids[0]
// 		sample_challenge_slice_start_index, sample_challenge_slice_end_index := storagechallenges.GetStorageChallengeSliceIndices(sample_total_data_length_in_bytes, sample_file_hash_string, sample_block_hash_string, sample_challenging_masternode_id)
// 		fmt.Println(fmt.Sprint(sample_challenge_slice_start_index))
// 		fmt.Println(fmt.Sprint(sample_challenge_slice_end_index))

// 		sample_file_path := storagechallenges.GetFilePathFromFileHash(sample_file_hash_string, samp_slice_of_file_hashes, sample_slice_of_file_paths)
// 		fmt.Println(sample_file_path)

// 		storagechallenges.AddMasternodesToDb(samp_slice_of_masternode_ids)

// 		storagechallenges.RemoveMasternodesFromDb(sample_top_n_closest_masternode_ids)

// 		storagechallenges.AddBlocksToDb(sample_slice_of_fake_block_hashes)

// 		sample_file_path_new := storagechallenges.GetFilePathFromFileHashUsingDb(sample_file_hash)
// 		fmt.Println(sample_file_path_new)

// 		rqsymbol_file_storage_data_folder_path := "D:\\je_golang_storage_challenge_code\\rqsymbol_files_stored_by_masternodes\\"
// 		storagechallenges.GenerateTestFoldersAndFiles(rqsymbol_file_storage_data_folder_path, number_of_storage_replicas)

// 		sample_masternode_id := samp_slice_of_masternode_ids[1]
// 		sample_file_hash2 := "9e63c16a4e6e5b29e614653660d059fa1d3285b2696c0452307db6f7426ff3f7"
// 		filepath_for_file_hash := storagechallenges.CheckForLocalFilepathForFileHashFunc(sample_masternode_id, sample_file_hash2, rqsymbol_file_storage_data_folder_path)
// 		fmt.Println(filepath_for_file_hash)

// 		challenges_per_masternode_per_block := int(math.Ceil(float64(len(samp_slice_of_masternode_ids)) / 2))
// 		number_of_challenge_replicas := 3
// 		sample_slice_of_message_ids := storagechallenges.GenerateStorageChallenges(sample_challenging_masternode_id, sample_block_hash_string, challenges_per_masternode_per_block, number_of_challenge_replicas, rqsymbol_file_storage_data_folder_path)
// 		fmt.Println(sample_slice_of_message_ids)

// 		storagechallenges.UpdateMasternodeStatsInDb()
// 		storagechallenges.UpdateBlockStatsInDb(sample_slice_of_fake_block_hashes)

// 		storagechallenges.RespondToStorageChallenges("jXamyHN3fe8Y8aLMdnVnMoGHJ8Y6nxVrLbwPpXwiarYxw5BROiKkZjSeVBqS6NrwluLrJXnVBhmg1hZLRsvtLq", rqsymbol_file_storage_data_folder_path, sample_block_hash_string)
// 	}

// 	cfg := &config.Config{}
// 	if err := cfg.Load(); err != nil {
// 		log.Panic("cannot load configuration file", err)
// 	}

// 	slice_of_new_raptorq_symbol_file_paths, _ := filepath.Glob(cfg.NewRqsymbolFileStorageDataFolderPath + "*")
// 	slice_of_raptorq_symbol_file_hashes := storagechallenges.GetRaptorqFileHashes(cfg.FolderPathContainingSampleRaptorqSymbolFiles)

// 	reset_simulation_state := true
// 	if reset_simulation_state {
// 		storagechallenges.ResetFolderState(cfg.RqsymbolFileStorageDataFolderPath)
// 		slice_of_file_paths_to_add_to_db, _ := filepath.Glob(cfg.FolderPathContainingSampleRaptorqSymbolFiles + "*")
// 		storagechallenges.AddFilesToDb(slice_of_file_paths_to_add_to_db)
// 		xor_distance_matrix := storagechallenges.ComputeMasternodeIdToFileHashXorDistanceMatrix(cfg.SliceOfPastelMasternodeIds, slice_of_raptorq_symbol_file_hashes)
// 		storagechallenges.AddXorDistanceMatrixToDb(cfg.SliceOfNewPastelMasternodeIds, slice_of_raptorq_symbol_file_hashes, xor_distance_matrix)
// 		storagechallenges.GenerateTestFoldersAndFiles(cfg.RqsymbolFileStorageDataFolderPath, cfg.NumberOfChallengeReplicas)
// 		storagechallenges.MakeDishonestMasternodesDeleteRandomFiles(cfg.SliceOfNewPastelMasternodeIds, cfg.RqsymbolFileStorageDataFolderPath)
// 	}

// 	/* 	xor_distances_slice := storagechallenges.GetXorDistancesFromDb()
// 	   	xor_distance_matrix := storagechallenges.TurnXorDistancesSliceIntoMatrix(xor_distances_slice, cfg.SliceOfNewPastelMasternodeIds, slice_of_raptorq_symbol_file_hashes)
// 	   	fmt.Println(xor_distance_matrix[0][0]) */
// 	xor_distance_matrix := storagechallenges.ComputeMasternodeIdToFileHashXorDistanceMatrix(cfg.SliceOfNewPastelMasternodeIds, slice_of_raptorq_symbol_file_hashes)

// 	initialize_database := true

// 	slice_of_block_hashes := make([]string, 0)
// 	if initialize_database {
// 		xor_distance_matrix := storagechallenges.ComputeMasternodeIdToFileHashXorDistanceMatrix(cfg.SliceOfNewPastelMasternodeIds, slice_of_raptorq_symbol_file_hashes)
// 		storagechallenges.AddXorDistanceMatrixToDb(cfg.SliceOfNewPastelMasternodeIds, slice_of_raptorq_symbol_file_hashes, xor_distance_matrix)
// 		number_of_blocks_to_make := 40
// 		slice_of_block_hashes = helper.GenerateFakeBlockHashes(number_of_blocks_to_make)
// 		fmt.Println("Adding files to database...")
// 		slice_of_file_paths_to_add_to_db, _ := filepath.Glob(cfg.FolderPathContainingSampleRaptorqSymbolFiles + "*")
// 		storagechallenges.AddFilesToDb(slice_of_file_paths_to_add_to_db)
// 		fmt.Println("Done!")
// 		storagechallenges.AddMasternodesToDb(cfg.SliceOfNewPastelMasternodeIds)
// 		storagechallenges.AddBlocksToDb(slice_of_block_hashes)
// 	}

// 	run_simulation := true

// 	if run_simulation {
// 		block_number := 1
// 		for _, current_block_hash := range slice_of_block_hashes {
// 			fmt.Println("\n\n_____________________________________________________________________________________________________________")
// 			fmt.Println("\n\nCurrent Block Number: " + fmt.Sprint(block_number) + " | Block Hash: " + current_block_hash)
// 			if ((block_number)%1 == 0) && (block_number <= 100) {
// 				n := 2
// 				k := 60
// 				storagechallenges.AddNIncrementalMasternodeIdsAndKIncrementalFiles(n, k, cfg.SliceOfNewPastelMasternodeIds, slice_of_new_raptorq_symbol_file_paths, xor_distance_matrix)
// 				storagechallenges.RedistributeFilesToMasternodes(cfg.RqsymbolFileStorageDataFolderPath, cfg.NumberOfChallengeReplicas)
// 				storagechallenges.MakeDishonestMasternodesDeleteRandomFiles(cfg.SliceOfNewPastelMasternodeIds, cfg.RqsymbolFileStorageDataFolderPath)
// 			}
// 			slice_of_masternode_ids, _ := storagechallenges.GetCurrentListsOfMasternodeIdsAndFileHashesFromDb()
// 			number_of_masternodes_to_issue_challenges_per_block := int(math.Ceil(float64(len(slice_of_masternode_ids)) / 3))
// 			challenges_per_masternode_per_block := int(math.Ceil(float64(len(slice_of_masternode_ids)) / 3))
// 			slice_of_challenging_masternode_ids_for_block := storagechallenges.GetNClosestMasternodeIdsToAGivenComparisonString(number_of_masternodes_to_issue_challenges_per_block, current_block_hash, slice_of_masternode_ids)
// 			storagechallenges.UpdateBlockStatsInDb(slice_of_block_hashes)
// 			for _, current_masternode_id := range slice_of_challenging_masternode_ids_for_block {
// 				slice_of_masternode_ids, _ := storagechallenges.GetCurrentListsOfMasternodeIdsAndFileHashesFromDb()
// 				_ = storagechallenges.GenerateStorageChallenges(current_masternode_id, current_block_hash, challenges_per_masternode_per_block, cfg.NumberOfChallengeReplicas, cfg.RqsymbolFileStorageDataFolderPath)
// 				for _, current_responding_masternode_id := range slice_of_masternode_ids {
// 					_ = storagechallenges.RespondToStorageChallenges(current_responding_masternode_id, cfg.RqsymbolFileStorageDataFolderPath, current_block_hash)
// 				}
// 				for _, current_verifying_masternode_id := range slice_of_masternode_ids {
// 					_ = storagechallenges.VerifyStorageChallengeResponses(current_verifying_masternode_id, cfg.RqsymbolFileStorageDataFolderPath, cfg.MaxSecondsToRespondToStorageChallenge)
// 				}
// 				storagechallenges.UpdateMasternodeStatsInDb()
// 				fmt.Println("\n\n_____________________________________________________________________________________________________________")
// 				time.Sleep(500 * time.Millisecond)
// 			}
// 			block_number += 1
// 		}
// 	}
// }
