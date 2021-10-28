package storagechallenges

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/karrick/godirwalk"
	"github.com/mkmik/argsort"
	"github.com/pastelnetwork/storage-challenges/external/storage"
	"github.com/pastelnetwork/storage-challenges/utils/file"
	"github.com/pastelnetwork/storage-challenges/utils/helper"
	"golang.org/x/crypto/sha3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetRaptorqFileHashes(path_to_raptorq_files string) []string {
	defer ExecutionDuration(TrackTime("GetRaptorqFileHashes"))
	d, _ := os.ReadDir(path_to_raptorq_files)
	number_of_files_in_folder := len(d)
	raptorq_file_hashes := make([]string, number_of_files_in_folder)
	matching_file_paths, _ := filepath.Glob(path_to_raptorq_files + "*")
	for idx, current_file_path := range matching_file_paths {
		raptorq_file_hashes[idx], _ = file.GetHashFromFilePath(current_file_path)
	}
	return raptorq_file_hashes
}

func TrackTime(msg string) (string, time.Time) {
	return msg, time.Now()
}

func ExecutionDuration(msg string, start time.Time) {
	log.Printf("%v: %v\n", msg, time.Since(start))
}

func ComputeMasternodeIdToFileHashXorDistanceMatrix(pastel_masternode_ids []string, raptorq_symbol_file_hashes []string) [][]uint64 {
	defer ExecutionDuration(TrackTime("ComputeMasternodeIdToFileHashXorDistanceMatrix"))
	fmt.Println("Generating XOR distance matrix...")
	xor_distance_matrix := make([][]uint64, len(pastel_masternode_ids))
	for ii := range xor_distance_matrix {
		xor_distance_matrix[ii] = make([]uint64, len(raptorq_symbol_file_hashes))
	}
	for idx1, current_masternode_id := range pastel_masternode_ids {
		for idx2, current_file_hash := range raptorq_symbol_file_hashes {
			xor_distance_matrix[idx1][idx2] = helper.ComputeXorDistanceBetweenTwoStrings(current_masternode_id, current_file_hash)
		}
	}
	return xor_distance_matrix
}

type XOR_Distance struct {
	gorm.Model
	Xor_distance_id string
	Masternode_id   string
	File_hash       string
	Xor_distance    uint64
}

func GetXorDistancesFromDb() []XOR_Distance {
	db := storage.GetStore().GetDB()
	var xor_distances_slice []XOR_Distance
	_ = db.Table("xor_distances").Select("*").Find(&xor_distances_slice)
	return xor_distances_slice
}

func TurnXorDistancesSliceIntoMatrix(xor_distances_slice []XOR_Distance, pastel_masternode_ids []string, raptorq_symbol_file_hashes []string) [][]uint64 {
	xor_distance_matrix := make([][]uint64, len(pastel_masternode_ids))
	for ii := range xor_distance_matrix {
		xor_distance_matrix[ii] = make([]uint64, len(raptorq_symbol_file_hashes))
	}
	column_count := 0
	row_count := 0
	for _, current_xor_distance_slice := range xor_distances_slice {
		if column_count >= len(raptorq_symbol_file_hashes) {
			column_count = 0
			row_count += 1
		}
		xor_distance_matrix[row_count][column_count] = current_xor_distance_slice.Xor_distance
		column_count += 1
	}
	return xor_distance_matrix
}

func FindMissingElementsOfAinB(a, b []string) []string {
	type void struct{}
	ma := make(map[string]void, len(a))
	diffs := []string{}
	for _, ka := range a {
		ma[ka] = void{}
	}
	for _, kb := range b {
		if _, ok := ma[kb]; !ok {
			diffs = append(diffs, kb)
		}
	}
	return diffs
}

func sliceContainsString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func AddXorDistanceMatrixToDb(pastel_masternode_ids []string, raptorq_symbol_file_hashes []string, xor_distance_matrix [][]uint64) error {
	defer ExecutionDuration(TrackTime("AddXorDistanceMatrixToDb"))
	existing_xor_distances_slice := GetXorDistancesFromDb()
	slice_of_existing_xor_distance_ids := make([]string, len(existing_xor_distances_slice))
	cnt1 := 0
	for _, current_xordistance := range existing_xor_distances_slice {
		slice_of_existing_xor_distance_ids[cnt1] = current_xordistance.Xor_distance_id
		cnt1 += 1
	}
	IncrementalXorDistanceSliceOfStructs := make([]XOR_Distance, 0)
	cnt := 0
	for idx1, current_masternode_id := range pastel_masternode_ids {
		for idx2, current_file_hash := range raptorq_symbol_file_hashes {
			current_xor_distance_id := helper.GetHashFromString(current_masternode_id + current_file_hash)
			if !sliceContainsString(slice_of_existing_xor_distance_ids, current_xor_distance_id) {
				current_xor_distance2 := xor_distance_matrix[idx1][idx2]
				if len(current_file_hash) > 0 {
					IncrementalXorDistanceSliceOfStructs = append(IncrementalXorDistanceSliceOfStructs, XOR_Distance{Xor_distance_id: current_xor_distance_id, Masternode_id: current_masternode_id, File_hash: current_file_hash, Xor_distance: current_xor_distance2})
					cnt += 1
				}
			}
		}
	}
	db := storage.GetStore().GetDB()
	tx := db.Begin()
	var err error
	defer func() {
		if err != nil {
			log.Println("ERROR: ", err)
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	err = tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "Xor_distance_id"}}, UpdateAll: true}).Create(&IncrementalXorDistanceSliceOfStructs).Error
	return err
}

func DetermineWhichMasternodesAreResponsibleForWhichFileHashes(number_of_challenge_replicas int) []XOR_Distance {
	db := storage.GetStore().GetDB()
	// TODO: improve subquery to use join instead, would be more faster, the cose to use this is too big
	// 	QUERY PLAN
	// |--CO-ROUTINE 2
	// |  |--CO-ROUTINE 4
	// |  |  |--SCAN TABLE xor_distances
	// |  |  `--USE TEMP B-TREE FOR ORDER BY
	// |  `--SCAN SUBQUERY 4
	// `--SCAN SUBQUERY 2
	//
	count_query_statement := "SELECT count(*) FROM (SELECT *, RANK() OVER (PARTITION BY file_hash ORDER BY xor_distance ASC) as rnk FROM (SELECT * FROM xor_distances)) WHERE rnk <= %v"
	count_query_statement_prepared := fmt.Sprintf(count_query_statement, number_of_challenge_replicas)
	var count_of_results int

	err := db.Raw(count_query_statement_prepared).Scan(&count_of_results).Error
	if err != nil {
		log.Fatal(err)
	}

	// TODO: improve same as above
	query_statement := "SELECT xor_distance_id, masternode_id, file_hash, xor_distance FROM (SELECT *, RANK() OVER (PARTITION BY file_hash ORDER BY xor_distance ASC) as rnk FROM (SELECT * FROM xor_distances)) WHERE rnk <= %v"
	query_statement_prepared := fmt.Sprintf(query_statement, number_of_challenge_replicas)
	rows, _ := db.Raw(query_statement_prepared).Rows()
	MasternodeToFileHashResponsibilityStructs := make([]XOR_Distance, 0)
	defer rows.Close()
	for rows.Next() {
		var xor_distance_id string
		var masternode_id string
		var file_hash string
		var xor_distance uint64
		rows.Scan(&xor_distance_id, &masternode_id, &file_hash, &xor_distance)
		var current_xor_distance XOR_Distance
		current_xor_distance.Xor_distance_id = xor_distance_id
		current_xor_distance.Masternode_id = masternode_id
		current_xor_distance.File_hash = file_hash
		current_xor_distance.Xor_distance = xor_distance
		if len(file_hash) == 64 {
			MasternodeToFileHashResponsibilityStructs = append(MasternodeToFileHashResponsibilityStructs, current_xor_distance)
		}
	}
	return MasternodeToFileHashResponsibilityStructs
}

func GetCurrentListsOfMasternodeIdsAndFileHashesFromDb() ([]string, []string) {
	db := storage.GetStore().GetDB()
	slice_of_masternode_ids := make([]string, 0)
	query_statement1 := "SELECT DISTINCT masternode_id FROM xor_distances"
	rows1, _ := db.Raw(query_statement1).Rows()
	defer rows1.Close()
	for rows1.Next() {
		var masternode_id string
		rows1.Scan(&masternode_id)
		slice_of_masternode_ids = append(slice_of_masternode_ids, masternode_id)
	}
	slice_of_file_hashes := make([]string, 0)
	query_statement2 := "SELECT DISTINCT file_hash FROM xor_distances"
	rows2, _ := db.Raw(query_statement2).Rows()
	defer rows2.Close()
	for rows2.Next() {
		var file_hash string
		rows2.Scan(&file_hash)
		slice_of_file_hashes = append(slice_of_file_hashes, file_hash)
	}
	return slice_of_masternode_ids, slice_of_file_hashes
}

func GetSliceOfFilePathsFromSliceOfFileHashes(slice_of_file_hashes []string) []string {
	db := storage.GetStore().GetDB()
	slice_of_file_paths := make([]string, 0)
	for _, current_file_hash := range slice_of_file_hashes {
		query_statement := "SELECT original_file_path FROM symbol_files WHERE file_hash = \"%v\""
		query_statement_prepared := fmt.Sprintf(query_statement, current_file_hash)
		rows, _ := db.Raw(query_statement_prepared).Rows()
		defer rows.Close()
		for rows.Next() {
			var current_file_path string
			rows.Scan(&current_file_path)
			if len(current_file_path) > 0 {
				slice_of_file_paths = append(slice_of_file_paths, current_file_path)
			}
		}
	}
	return slice_of_file_paths
}

func GetNClosestMasternodeIdsToAGivenFileHashUsingDb(n int, file_hash_string string) []string {
	db := storage.GetStore().GetDB()
	slice_of_top_n_closest_masternode_ids := make([]string, 0)
	query_statement := "SELECT masternode_id FROM xor_distances WHERE file_hash=\"%s\" ORDER BY xor_distance ASC LIMIT %d"
	query_statement_prepared := fmt.Sprintf(query_statement, file_hash_string, n)
	rows, _ := db.Raw(query_statement_prepared).Rows()
	defer rows.Close()
	for rows.Next() {
		var masternode_id string
		rows.Scan(&masternode_id)
		if len(masternode_id) > 0 {
			slice_of_top_n_closest_masternode_ids = append(slice_of_top_n_closest_masternode_ids, masternode_id)
		}
	}
	return slice_of_top_n_closest_masternode_ids
}

func GetNClosestMasternodeIdsToAGivenComparisonString(n int, comparison_string string, slice_of_pastel_masternode_ids []string) []string {
	slice_of_xor_distances := make([]uint64, len(slice_of_pastel_masternode_ids))
	XORdistance_to_masternodeID_map := make(map[uint64]string)
	for idx, current_masternode_id := range slice_of_pastel_masternode_ids {
		current_xor_distance := helper.ComputeXorDistanceBetweenTwoStrings(current_masternode_id, comparison_string)
		slice_of_xor_distances[idx] = current_xor_distance
		XORdistance_to_masternodeID_map[current_xor_distance] = current_masternode_id
	}
	sort.Slice(slice_of_xor_distances, func(i, j int) bool { return slice_of_xor_distances[i] < slice_of_xor_distances[j] })
	slice_of_top_n_closest_masternode_ids := make([]string, n)
	for ii, current_xor_distance := range slice_of_xor_distances {
		if ii < n {
			slice_of_top_n_closest_masternode_ids[ii] = XORdistance_to_masternodeID_map[current_xor_distance]
		}
	}
	return slice_of_top_n_closest_masternode_ids
}

func GetNClosestFileHashesToAGivenComparisonString(n int, comparison_string string, slice_of_file_hashes []string) []string {
	slice_of_xor_distances := make([]uint64, len(slice_of_file_hashes))
	XORdistance_to_fileHash_map := make(map[uint64]string)
	for idx, current_file_hash := range slice_of_file_hashes {
		current_xor_distance := helper.ComputeXorDistanceBetweenTwoStrings(current_file_hash, comparison_string)
		slice_of_xor_distances[idx] = current_xor_distance
		XORdistance_to_fileHash_map[current_xor_distance] = current_file_hash
	}
	sort.Slice(slice_of_xor_distances, func(i, j int) bool { return slice_of_xor_distances[i] < slice_of_xor_distances[j] })
	slice_of_top_n_closest_file_hashes := make([]string, n)
	for ii, current_xor_distance := range slice_of_xor_distances {
		if ii < n {
			slice_of_top_n_closest_file_hashes[ii] = XORdistance_to_fileHash_map[current_xor_distance]
		}
	}
	return slice_of_top_n_closest_file_hashes
}

func MinMax(array []int) (int, int) {
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

func GetFilePathsFromFolder(path_to_raptorq_files string) []string {
	raptorq_file_paths, _ := filepath.Glob(path_to_raptorq_files + "*")
	return raptorq_file_paths
}

func GetFilePathFromFileHash(file_hash_string string, slice_of_raptorq_symbol_file_hashes []string, slice_of_raptorq_symbol_file_paths []string) string {
	file_hash_to_path_map := make(map[string]string)
	for idx, current_file_hash := range slice_of_raptorq_symbol_file_hashes {
		current_file_path := slice_of_raptorq_symbol_file_paths[idx]
		file_hash_to_path_map[current_file_hash] = current_file_path
	}
	file_path := file_hash_to_path_map[file_hash_string]
	return file_path
}

func ComputeHashofFileSlice(file_data []byte, challenge_slice_start_index int, challenge_slice_end_index int) string {
	challenge_data_slice := file_data[challenge_slice_start_index:challenge_slice_end_index]
	algorithm := sha3.New256()
	algorithm.Write(challenge_data_slice)
	hash_of_data_slice := hex.EncodeToString(algorithm.Sum(nil))
	return hash_of_data_slice
}

type SymbolFiles struct {
	gorm.Model
	File_hash                 string
	File_length_in_bytes      uint
	Total_challenges_for_file uint
	Original_file_path        string
}

func AddFilesToDb(slice_of_input_file_paths []string) {
	defer ExecutionDuration(TrackTime("AddFilesToDb"))
	db := storage.GetStore().GetDB()
	tx := db.Begin()
	var err error
	defer func() {
		if err != nil {
			log.Println("ERROR: ", err)
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	var listSymbolfiles = make([]SymbolFiles, 0)
	// use batch query instead: INSERT INTO symbol_files (col1, col2, col3) VALUES (?, ?, ?), (?, ?, ?) on conflict do update set col1=excluded.col1, col2=excluded.col2, col3=excluded.col3;
	for _, current_input_file_path := range slice_of_input_file_paths {
		current_file_hash, _ := file.GetHashFromFilePath(current_input_file_path)
		var current_symbolfile_struct SymbolFiles
		current_symbolfile_struct.File_hash = current_file_hash
		fi, _ := os.Stat(current_input_file_path)
		current_file_size := fi.Size()
		current_symbolfile_struct.File_length_in_bytes = uint(current_file_size)
		current_symbolfile_struct.Total_challenges_for_file = uint(0)
		current_symbolfile_struct.Original_file_path = current_input_file_path
		listSymbolfiles = append(listSymbolfiles, current_symbolfile_struct)
	}
	err = tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "File_hash"}}, UpdateAll: true}).Create(&listSymbolfiles).Error
}

type Masternodes struct {
	gorm.Model
	Masternode_id                         string
	Masternode_ip_address                 string
	Total_challenges_issued               uint
	Total_challenges_responded_to         uint
	Total_challenges_correct              uint
	Total_challenges_incorrect            uint
	Total_challenges_correct_but_too_slow uint
	Total_challenges_never_responded_to   uint
	Challenge_response_success_rate_pct   float32
}

func AddMasternodesToDb(slice_of_pastel_masternode_ids []string) {
	db := storage.GetStore().GetDB()
	var err error
	tx := db.Begin()
	defer func() {
		if err != nil {
			log.Println("ERROR: ", err)
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	var list_masternodes = make([]Masternodes, 0)
	for _, current_masternode_id := range slice_of_pastel_masternode_ids {
		var current_masternode_struct Masternodes
		current_masternode_struct.Masternode_id = current_masternode_id
		current_masternode_struct.Masternode_ip_address = "127.0.0.1"
		list_masternodes = append(list_masternodes, current_masternode_struct)
	}

	err = tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "Masternode_id"}}, UpdateAll: true}).Create(&list_masternodes).Error
}

func RemoveMasternodesFromDb(slice_of_pastel_masternode_ids []string) {
	db := storage.GetStore().GetDB()
	var err error
	tx := db.Begin()
	defer func() {
		if err != nil {
			log.Println("ERROR: ", err)
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	err = tx.Unscoped().Delete(&Masternodes{}, "Masternode_id in (?) ", slice_of_pastel_masternode_ids).Error
}

type PastelBlocks struct {
	gorm.Model
	Block_hash                            string
	Block_number                          uint
	Total_challenges_issued               uint
	Total_challenges_responded_to         uint
	Total_challenges_correct              uint
	Total_challenges_incorrect            uint
	Total_challenges_correct_but_too_slow uint
	Total_challenges_never_responded_to   uint
	Challenge_response_success_rate_pct   float32
}

func AddBlocksToDb(slice_of_block_hashes []string) {
	db := storage.GetStore().GetDB()
	var err error
	tx := db.Begin()
	defer func() {
		if err != nil {
			log.Println("ERROR: ", err)
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	var list_pastel_blocks = make([]PastelBlocks, 0)
	for idx, current_block_hash := range slice_of_block_hashes {
		var current_pastel_block PastelBlocks
		current_pastel_block.Block_hash = current_block_hash
		current_pastel_block.Block_number = uint(idx)
		list_pastel_blocks = append(list_pastel_blocks, current_pastel_block)
	}
	err = tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "Block_hash"}}, UpdateAll: true}).Create(&list_pastel_blocks).Error
}

func RemoveGlob(path string) (err error) {
	contents, err := filepath.Glob(path)
	if err != nil {
		return
	}
	for _, item := range contents {
		err = os.RemoveAll(item)
		if err != nil {
			return
		}
	}
	return
}

func pruneEmptyDirectories(osDirname string) (int, error) {
	var count int
	err := godirwalk.Walk(osDirname, &godirwalk.Options{
		Unsorted: true,
		Callback: func(_ string, _ *godirwalk.Dirent) error {
			return nil
		},
		PostChildrenCallback: func(osPathname string, _ *godirwalk.Dirent) error {
			s, err := godirwalk.NewScanner(osPathname)
			if err != nil {
				return err
			}
			hasAtLeastOneChild := s.Scan()
			if err := s.Err(); err != nil {
				return err
			}
			if hasAtLeastOneChild {
				return nil
			}
			if osPathname == osDirname {
				return nil
			}
			err = os.Remove(osPathname)
			if err == nil {
				count++
			}
			return err
		},
	})
	return count, err
}

func ResetFolderState(rqsymbol_file_storage_data_folder_path string) {
	defer ExecutionDuration(TrackTime("ResetFolderState"))
	fmt.Println("Resetting rqsymbol storage folders (and files) by deleting them!")
	RemoveGlob(rqsymbol_file_storage_data_folder_path + "*/*")
	pruneEmptyDirectories(rqsymbol_file_storage_data_folder_path)
}

func ensureDir(dirName string) error {
	err := os.Mkdir(dirName, os.ModeDir)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		info, err := os.Stat(dirName)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return errors.New("path exists but is not a directory")
		}
		return nil
	}
	return err
}

func GetFilePathFromFileHashUsingDb(file_hash_string string) string {
	db := storage.GetStore().GetDB()
	var file_path string
	row := db.Table("symbol_files").Where("file_hash = ?", file_hash_string).Select("original_file_path").Row()
	row.Scan(&file_path)
	return file_path
}

func GetSliceOfFilePathsMasternodeIsResponsibleFor(masternode_id string, MasternodeToFileHashResponsibilityStructs []XOR_Distance) []string {
	slice_of_file_paths_masternode_is_responsible_for := make([]string, 0)
	for _, current_xor_distance_struct := range MasternodeToFileHashResponsibilityStructs {
		if current_xor_distance_struct.Masternode_id == masternode_id {
			current_file_hash := current_xor_distance_struct.File_hash
			current_file_path := GetFilePathFromFileHashUsingDb(current_file_hash)
			slice_of_file_paths_masternode_is_responsible_for = append(slice_of_file_paths_masternode_is_responsible_for, current_file_path)
		}
	}
	return slice_of_file_paths_masternode_is_responsible_for
}

func GetSliceOfFileHashesMasternodeIsResponsibleFor(masternode_id string, MasternodeToFileHashResponsibilityStructs []XOR_Distance) []string {
	slice_of_file_hashes_masternode_is_responsible_for := make([]string, 0)
	for _, current_xor_distance_struct := range MasternodeToFileHashResponsibilityStructs {
		if current_xor_distance_struct.Masternode_id == masternode_id {
			current_file_hash := current_xor_distance_struct.File_hash
			slice_of_file_hashes_masternode_is_responsible_for = append(slice_of_file_hashes_masternode_is_responsible_for, current_file_hash)
		}
	}
	return slice_of_file_hashes_masternode_is_responsible_for
}

func resetReadOnlyFlagAll(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return os.Chmod(path, 0666)
	}
	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fd.Close()
	names, _ := fd.Readdirnames(-1)
	for _, name := range names {
		resetReadOnlyFlagAll(path + string(filepath.Separator) + name)
	}
	return nil
}

func CopyFile(src, dst string) error {
	BUFFERSIZE := int64(1000)
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()
	_, err = os.Stat(dst)
	if err == nil {
		return fmt.Errorf("file %s already exists", dst)
	}
	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	if err != nil {
		panic(err)
	}
	buf := make([]byte, BUFFERSIZE)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}

func CopyFilesForMasternode(current_masternode_id string, rqsymbol_file_storage_data_folder_path string, MasternodeToFileHashResponsibilityStructs []XOR_Distance) {
	slice_of_file_hashes_masternode_is_responsible_for := GetSliceOfFileHashesMasternodeIsResponsibleFor(current_masternode_id, MasternodeToFileHashResponsibilityStructs)
	slice_of_file_paths_masternode_is_responsible_for := GetSliceOfFilePathsMasternodeIsResponsibleFor(current_masternode_id, MasternodeToFileHashResponsibilityStructs)
	for idx, current_file_path := range slice_of_file_paths_masternode_is_responsible_for {
		current_file_hash := slice_of_file_hashes_masternode_is_responsible_for[idx]
		current_masternode_folder_path := rqsymbol_file_storage_data_folder_path + current_masternode_id
		renamed_destination_file_path := current_masternode_folder_path + "/" + current_file_hash + ".rqs"
		CopyFile(current_file_path, renamed_destination_file_path)
		resetReadOnlyFlagAll(renamed_destination_file_path)
	}
}

func GenerateTestFoldersAndFiles(rqsymbol_file_storage_data_folder_path string, number_of_challenge_replicas int) {
	defer ExecutionDuration(TrackTime("GenerateTestFoldersAndFiles"))
	slice_of_masternode_ids, slice_of_file_hashes := GetCurrentListsOfMasternodeIdsAndFileHashesFromDb()
	if _, err := os.Stat(rqsymbol_file_storage_data_folder_path); os.IsNotExist(err) {
		ensureDir(rqsymbol_file_storage_data_folder_path)
	}
	for _, current_masternode_id := range slice_of_masternode_ids {
		current_masternode_folder_path := rqsymbol_file_storage_data_folder_path + current_masternode_id
		if _, err := os.Stat(current_masternode_folder_path); os.IsNotExist(err) {
			ensureDir(current_masternode_folder_path)
		}
	}
	fmt.Println("Assigning " + fmt.Sprint(len(slice_of_file_hashes)) + " files to " + fmt.Sprint(len(slice_of_masternode_ids)) + " different masternodes...")
	MasternodeToFileHashResponsibilityStructs := DetermineWhichMasternodesAreResponsibleForWhichFileHashes(number_of_challenge_replicas)
	for _, current_masternode_id := range slice_of_masternode_ids {
		CopyFilesForMasternode(current_masternode_id, rqsymbol_file_storage_data_folder_path, MasternodeToFileHashResponsibilityStructs)
	}
	resetReadOnlyFlagAll(rqsymbol_file_storage_data_folder_path)
}

func RedistributeFilesToMasternodes(rqsymbol_file_storage_data_folder_path string, number_of_challenge_replicas int) {
	defer ExecutionDuration(TrackTime("RedistributeFilesToMasternodes"))
	slice_of_masternode_ids, slice_of_file_hashes := GetCurrentListsOfMasternodeIdsAndFileHashesFromDb()
	if _, err := os.Stat(rqsymbol_file_storage_data_folder_path); os.IsNotExist(err) {
		ensureDir(rqsymbol_file_storage_data_folder_path)
	}
	for _, current_masternode_id := range slice_of_masternode_ids {
		current_masternode_folder_path := rqsymbol_file_storage_data_folder_path + current_masternode_id
		if _, err := os.Stat(current_masternode_folder_path); os.IsNotExist(err) {
			ensureDir(current_masternode_folder_path)
		}
	}
	fmt.Println("Redistributing " + fmt.Sprint(len(slice_of_file_hashes)) + " files to " + fmt.Sprint(len(slice_of_masternode_ids)) + " different masternodes...")
	MasternodeToFileHashResponsibilityStructs := DetermineWhichMasternodesAreResponsibleForWhichFileHashes(number_of_challenge_replicas)
	for _, current_masternode_id := range slice_of_masternode_ids {
		slice_of_file_hashes_masternode_is_responsible_for := GetSliceOfFileHashesMasternodeIsResponsibleFor(current_masternode_id, MasternodeToFileHashResponsibilityStructs)
		current_masternode_folder_path := rqsymbol_file_storage_data_folder_path + current_masternode_id
		slice_of_file_hashes_currently_stored_by_masternode := make([]string, 0)
		matching_file_paths, _ := filepath.Glob(current_masternode_folder_path + "/" + "*")
		for _, current_file_path := range matching_file_paths {
			_, current_file_name := filepath.Split(current_file_path)
			current_file_hash := strings.Replace(current_file_name, ".rqs", "", -1)
			slice_of_file_hashes_currently_stored_by_masternode = append(slice_of_file_hashes_currently_stored_by_masternode, current_file_hash)
		}
		slice_of_new_file_hashes_for_masternode_to_store := FindMissingElementsOfAinB(slice_of_file_hashes_masternode_is_responsible_for, slice_of_file_hashes_currently_stored_by_masternode)
		slice_of_file_hashes_masternode_is_storing_but_no_longer_has_to := FindMissingElementsOfAinB(slice_of_file_hashes_currently_stored_by_masternode, slice_of_file_hashes_masternode_is_responsible_for)
		fmt.Println("Masternode " + current_masternode_id + " is required to store an additional " + fmt.Sprint(len(slice_of_new_file_hashes_for_masternode_to_store)) + " files, and is no longer responsible for " + fmt.Sprint(len(slice_of_file_hashes_masternode_is_storing_but_no_longer_has_to)))
		for _, current_file_hash := range slice_of_new_file_hashes_for_masternode_to_store {
			current_file_path := GetFilePathFromFileHashUsingDb(current_file_hash)
			current_masternode_folder_path := rqsymbol_file_storage_data_folder_path + current_masternode_id
			renamed_destination_file_path1 := current_masternode_folder_path + "/" + current_file_hash + ".rqs"
			CopyFile(current_file_path, renamed_destination_file_path1)
		}
		for _, current_file_hash2 := range slice_of_file_hashes_masternode_is_storing_but_no_longer_has_to {
			current_masternode_folder_path2 := rqsymbol_file_storage_data_folder_path + current_masternode_id
			renamed_destination_file_path2 := current_masternode_folder_path2 + "/" + current_file_hash2 + ".rqs"
			resetReadOnlyFlagAll(renamed_destination_file_path2)
			RemoveGlob(renamed_destination_file_path2)
		}
	}
}

func CheckForLocalFilepathForFileHashFunc(masternode_id string, file_hash string, rqsymbol_file_storage_data_folder_path string) string {
	masternode_storage_path := rqsymbol_file_storage_data_folder_path + masternode_id + "/"
	masternode_storage_path_glob_matches, _ := filepath.Glob(masternode_storage_path + file_hash + ".rqs")
	filepath_for_file_hash := ""
	if len(masternode_storage_path_glob_matches) > 0 {
		filepath_for_file_hash = masternode_storage_path_glob_matches[0]
	}
	return filepath_for_file_hash
}

type ChallengeMessages struct {
	gorm.Model
	Message_id                      string
	Message_type                    string
	Challenge_status                string
	Datetime_challenge_sent         string
	Datetime_challenge_responded_to string
	Datetime_challenge_verified     string
	Block_hash_when_challenge_sent  string
	Challenging_masternode_id       string
	Responding_masternode_id        string
	File_hash_to_challenge          string
	Challenge_slice_start_index     uint64
	Challenge_slice_end_index       uint64
	Challenge_slice_correct_hash    string
	Challenge_response_hash         string
	Challenge_id                    string
}

type Challenges struct {
	gorm.Model
	Challenge_id                       string
	Challenge_status                   string
	Datetime_challenge_sent            string
	Datetime_challenge_responded_to    string
	Datetime_challenge_verified        string
	Block_hash_when_challenge_sent     string
	Challenge_response_time_in_seconds float64
	Challenging_masternode_id          string
	Responding_masternode_id           string
	File_hash_to_challenge             string
	Challenge_slice_start_index        uint64
	Challenge_slice_end_index          uint64
	Challenge_slice_correct_hash       string
	Challenge_response_hash            string
}

func UpdateDbWithMessage(storage_challenge_message ChallengeMessages) {
	db := storage.GetStore().GetDB()
	var err error
	tx := db.Begin()
	defer func() {
		if err != nil {
			log.Println("ERROR: ", err)
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	err = tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "Message_id"}}, UpdateAll: true}).Create(&storage_challenge_message).Error
	if err != nil {
		return
	}
	challenge_id_input_data := storage_challenge_message.Challenging_masternode_id + storage_challenge_message.Responding_masternode_id + storage_challenge_message.File_hash_to_challenge + fmt.Sprint(storage_challenge_message.Challenge_slice_start_index) + fmt.Sprint(storage_challenge_message.Challenge_slice_end_index) + storage_challenge_message.Datetime_challenge_sent
	challenge_id := helper.GetHashFromString(challenge_id_input_data)
	var storage_challenge Challenges
	storage_challenge.Challenge_id = challenge_id
	storage_challenge.Challenge_status = storage_challenge_message.Challenge_status
	storage_challenge.Datetime_challenge_sent = storage_challenge_message.Datetime_challenge_sent
	storage_challenge.Datetime_challenge_responded_to = storage_challenge_message.Datetime_challenge_responded_to
	storage_challenge.Datetime_challenge_verified = storage_challenge_message.Datetime_challenge_verified
	storage_challenge.Block_hash_when_challenge_sent = storage_challenge_message.Block_hash_when_challenge_sent
	storage_challenge.Challenging_masternode_id = storage_challenge_message.Challenging_masternode_id
	storage_challenge.Responding_masternode_id = storage_challenge_message.Responding_masternode_id
	storage_challenge.File_hash_to_challenge = storage_challenge_message.File_hash_to_challenge
	storage_challenge.Challenge_slice_start_index = storage_challenge_message.Challenge_slice_end_index
	storage_challenge.Challenge_slice_end_index = storage_challenge_message.Challenge_slice_end_index
	storage_challenge.Challenge_slice_correct_hash = storage_challenge_message.Challenge_slice_correct_hash
	err = tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "Challenge_id"}}, UpdateAll: true}).Create(&storage_challenge).Error
}

func GetStorageChallengeSliceIndices(total_data_length_in_bytes uint64, file_hash_string string, block_hash_string string, challenging_masternode_id string) (int, int) {
	block_hash_string_as_int, _ := strconv.ParseInt(block_hash_string, 16, 64)
	block_hash_string_as_int_str := fmt.Sprint(block_hash_string_as_int)
	step_size_for_indices_str := block_hash_string_as_int_str[len(block_hash_string_as_int_str)-1:] + block_hash_string_as_int_str[0:1]
	step_size_for_indices, _ := strconv.ParseUint(step_size_for_indices_str, 10, 32)
	step_size_for_indices_as_int := int(step_size_for_indices)
	comparison_string := block_hash_string + file_hash_string + challenging_masternode_id
	slice_of_xor_distances_of_indices_to_block_hash := make([]uint64, 0)
	slice_of_indices_with_step_size := make([]int, 0)
	total_data_length_in_bytes_as_int := int(total_data_length_in_bytes)
	for j := 0; j <= total_data_length_in_bytes_as_int; j += step_size_for_indices_as_int {
		j_as_string := fmt.Sprintf("%d", j)
		current_xor_distance := helper.ComputeXorDistanceBetweenTwoStrings(j_as_string, comparison_string)
		slice_of_xor_distances_of_indices_to_block_hash = append(slice_of_xor_distances_of_indices_to_block_hash, current_xor_distance)
		slice_of_indices_with_step_size = append(slice_of_indices_with_step_size, j)
	}
	slice_of_sorted_indices := argsort.SortSlice(slice_of_xor_distances_of_indices_to_block_hash, func(i, j int) bool {
		return slice_of_xor_distances_of_indices_to_block_hash[i] < slice_of_xor_distances_of_indices_to_block_hash[j]
	})
	slice_of_sorted_indices_with_step_size := make([]int, 0)
	for _, current_sorted_index := range slice_of_sorted_indices {
		slice_of_sorted_indices_with_step_size = append(slice_of_sorted_indices_with_step_size, slice_of_indices_with_step_size[current_sorted_index])
	}
	first_two_sorted_indices := slice_of_sorted_indices_with_step_size[0:2]
	challenge_slice_start_index, challenge_slice_end_index := MinMax(first_two_sorted_indices)
	return challenge_slice_start_index, challenge_slice_end_index
}

func GenerateStorageChallenges(challenging_masternode_id string, current_block_hash string, challenges_per_masternode_per_block int, number_of_challenge_replicas int, rqsymbol_file_storage_data_folder_path string) []string {
	defer ExecutionDuration(TrackTime("GenerateStorageChallenges"))
	slice_of_message_ids := make([]string, 0)
	_, slice_of_file_hashes := GetCurrentListsOfMasternodeIdsAndFileHashesFromDb()
	slice_to_check_if_file_contained_by_local_masternode := make([]bool, 0)
	for _, current_file_hash := range slice_of_file_hashes {
		filepath_for_file_hash := CheckForLocalFilepathForFileHashFunc(challenging_masternode_id, current_file_hash, rqsymbol_file_storage_data_folder_path)
		if len(filepath_for_file_hash) > 0 {
			slice_to_check_if_file_contained_by_local_masternode = append(slice_to_check_if_file_contained_by_local_masternode, true)
		} else {
			slice_to_check_if_file_contained_by_local_masternode = append(slice_to_check_if_file_contained_by_local_masternode, false)
		}
	}
	slice_of_file_hashes_stored_by_challenger := make([]string, 0)
	for idx, current_file_contained_by_local_mn := range slice_to_check_if_file_contained_by_local_masternode {
		if current_file_contained_by_local_mn {
			slice_of_file_hashes_stored_by_challenger = append(slice_of_file_hashes_stored_by_challenger, slice_of_file_hashes[idx])
		}
	}
	comparison_string_for_file_hash_selection := current_block_hash + challenging_masternode_id
	slice_of_file_hashes_to_challenge := GetNClosestFileHashesToAGivenComparisonString(challenges_per_masternode_per_block, comparison_string_for_file_hash_selection, slice_of_file_hashes_stored_by_challenger)
	slice_of_masternodes_to_challenge := make([]string, len(slice_of_file_hashes_to_challenge))
	fmt.Println("Challenging Masternode " + challenging_masternode_id + " is now selecting file hashes to challenge this block, and then for each one, selecting which Masternode to challenge...")
	for idx1, current_file_hash_to_challenge := range slice_of_file_hashes_to_challenge {
		slice_of_masternodes_storing_file_hash := GetNClosestMasternodeIdsToAGivenFileHashUsingDb(number_of_challenge_replicas, current_file_hash_to_challenge)
		slice_of_masternodes_storing_file_hash_excluding_challenger := make([]string, 0)
		for idx, current_mastetnode_id := range slice_of_masternodes_storing_file_hash {
			if current_mastetnode_id == challenging_masternode_id {
				slice_of_masternodes_storing_file_hash_excluding_challenger = append(slice_of_masternodes_storing_file_hash[:idx], slice_of_masternodes_storing_file_hash[idx+1:]...)
			}
		}
		comparison_string_for_masternode_selection := current_block_hash + current_file_hash_to_challenge + challenging_masternode_id + helper.GetHashFromString(fmt.Sprint(idx1))
		responding_masternode_id := GetNClosestMasternodeIdsToAGivenComparisonString(1, comparison_string_for_masternode_selection, slice_of_masternodes_storing_file_hash_excluding_challenger)
		slice_of_masternodes_to_challenge[idx1] = responding_masternode_id[0]
	}
	for idx2, current_file_hash_to_challenge := range slice_of_file_hashes_to_challenge {
		filepath_for_challenge_file_hash := CheckForLocalFilepathForFileHashFunc(challenging_masternode_id, current_file_hash_to_challenge, rqsymbol_file_storage_data_folder_path)
		if len(filepath_for_challenge_file_hash) > 0 {
			challenge_file_data, _ := file.ReadFileIntoMemory(filepath_for_challenge_file_hash)
			challenge_data_size := uint64(len(challenge_file_data))
			if challenge_data_size > 0 {
				responding_masternode_id := slice_of_masternodes_to_challenge[idx2]
				challenge_status := "Pending"
				message_type := "storage_challenge_issuance_message"
				challenge_slice_start_index, challenge_slice_end_index := GetStorageChallengeSliceIndices(challenge_data_size, current_file_hash_to_challenge, current_block_hash, challenging_masternode_id)
				datetime_challenge_sent := time.Now().Format(time.RFC3339)
				message_id_input_data := challenging_masternode_id + responding_masternode_id + current_file_hash_to_challenge + challenge_status + message_type + current_block_hash
				message_id := helper.GetHashFromString(message_id_input_data)
				slice_of_message_ids = append(slice_of_message_ids, message_id)
				challenge_id_input_data := challenging_masternode_id + responding_masternode_id + current_file_hash_to_challenge + fmt.Sprint(challenge_slice_start_index) + fmt.Sprint(challenge_slice_end_index) + datetime_challenge_sent
				challenge_id := helper.GetHashFromString(challenge_id_input_data)
				var storage_challenge_message ChallengeMessages
				storage_challenge_message.Message_id = message_id
				storage_challenge_message.Message_type = message_type
				storage_challenge_message.Challenge_status = challenge_status
				storage_challenge_message.Datetime_challenge_sent = datetime_challenge_sent
				storage_challenge_message.Datetime_challenge_responded_to = ""
				storage_challenge_message.Datetime_challenge_verified = ""
				storage_challenge_message.Block_hash_when_challenge_sent = current_block_hash
				storage_challenge_message.Challenging_masternode_id = challenging_masternode_id
				storage_challenge_message.Responding_masternode_id = responding_masternode_id
				storage_challenge_message.File_hash_to_challenge = current_file_hash_to_challenge
				storage_challenge_message.Challenge_slice_start_index = uint64(challenge_slice_start_index)
				storage_challenge_message.Challenge_slice_end_index = uint64(challenge_slice_end_index)
				storage_challenge_message.Challenge_slice_correct_hash = ""
				storage_challenge_message.Challenge_response_hash = ""
				storage_challenge_message.Challenge_id = challenge_id
				UpdateDbWithMessage(storage_challenge_message)
				fmt.Println("\nMasternode " + challenging_masternode_id + " issued a storage challenge to Masternode " + responding_masternode_id + " for file hash " + current_file_hash_to_challenge + " (start index: " + fmt.Sprint(challenge_slice_start_index) + "; end index: " + fmt.Sprint(challenge_slice_end_index) + ")")
			} else {
				fmt.Println("\nMasternode " + challenging_masternode_id + " encountered an invalid file while attempting to generate a storage challenge for file hash " + current_file_hash_to_challenge)
			}
		} else {
			fmt.Println("\nMasternode " + challenging_masternode_id + " encountered an error generating storage challenges")
		}
	}
	return slice_of_message_ids
}

func UpdateMasternodeStatsInDb() {
	slice_of_masternode_ids, _ := GetCurrentListsOfMasternodeIdsAndFileHashesFromDb()
	total_challenges_issued := 0
	total_challenges_responded_to := 0
	total_challenges_correct := 0
	total_challenges_incorrect := 0
	total_challenges_correct_but_too_slow := 0
	total_challenges_never_responded_to := 0
	challenge_response_success_rate_pct := float32(0.0)

	db := storage.GetStore().GetDB()
	var err error
	tx := db.Begin()
	defer func() {
		if err != nil {
			log.Println("ERROR: ", err)
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	for _, current_masternode_id := range slice_of_masternode_ids {
		var challenges_issued []ChallengeMessages
		err = tx.Where("responding_masternode_id = ? AND message_type = ?", current_masternode_id, "storage_challenge_issuance_message").Find(&challenges_issued).Error
		if err != nil {
			return
		}
		total_challenges_issued = len(challenges_issued)
		var challenges_responded_to []ChallengeMessages
		err = tx.Where("responding_masternode_id = ? AND challenge_status = ?", current_masternode_id, "Responded").Find(&challenges_responded_to).Error
		if err != nil {
			return
		}
		total_challenges_responded_to = len(challenges_responded_to)
		var challenges_correct []ChallengeMessages
		err = tx.Where("responding_masternode_id = ? AND challenge_status = ?", current_masternode_id, "Successful response").Find(&challenges_correct).Error
		if err != nil {
			return
		}
		total_challenges_correct = len(challenges_correct)
		var challenges_incorrect []ChallengeMessages
		err = tx.Where("responding_masternode_id = ? AND challenge_status = ?", current_masternode_id, "Failed because of incorrect response").Find(&challenges_incorrect).Error
		if err != nil {
			return
		}
		total_challenges_incorrect = len(challenges_incorrect)
		var challenges_correct_but_too_slow []ChallengeMessages
		err = tx.Where("responding_masternode_id = ? AND challenge_status = ?", current_masternode_id, "Failed--Correct but response was too slow").Find(&challenges_correct_but_too_slow).Error
		if err != nil {
			return
		}
		total_challenges_correct_but_too_slow = len(challenges_correct_but_too_slow)
		var challenges_never_responded_to []ChallengeMessages
		err = tx.Where("responding_masternode_id = ? AND challenge_status = ?", current_masternode_id, "Failed because response never arrived").Find(&challenges_never_responded_to).Error
		if err != nil {
			return
		}
		total_challenges_never_responded_to = len(challenges_never_responded_to)
		if total_challenges_issued > 0 {
			challenge_response_success_rate_pct = float32(total_challenges_correct) / float32(total_challenges_issued)
		} else {
			challenge_response_success_rate_pct = float32(1.0)
		}
		var current_masternode_update Masternodes
		current_masternode_update.Masternode_id = current_masternode_id
		current_masternode_update.Total_challenges_issued = uint(total_challenges_issued)
		current_masternode_update.Total_challenges_responded_to = uint(total_challenges_responded_to)
		current_masternode_update.Total_challenges_incorrect = uint(total_challenges_incorrect)
		current_masternode_update.Total_challenges_correct_but_too_slow = uint(total_challenges_correct_but_too_slow)
		current_masternode_update.Total_challenges_never_responded_to = uint(total_challenges_never_responded_to)
		current_masternode_update.Challenge_response_success_rate_pct = challenge_response_success_rate_pct
		err = tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "Masternode_id"}}, UpdateAll: true}).Create(&current_masternode_update).Error
	}
}

func UpdateBlockStatsInDb(slice_of_block_hashes []string) {
	defer ExecutionDuration(TrackTime("UpdateBlockStatsInDb"))
	total_challenges_issued := 0
	total_challenges_responded_to := 0
	total_challenges_correct := 0
	total_challenges_incorrect := 0
	total_challenges_correct_but_too_slow := 0
	total_challenges_never_responded_to := 0
	challenge_response_success_rate_pct := float32(0.0)
	db := storage.GetStore().GetDB()
	var err error
	tx := db.Begin()
	defer func() {
		if err != nil {
			log.Println("ERROR: ", err)
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	for _, current_block_hash := range slice_of_block_hashes {
		var challenges_issued []ChallengeMessages
		err = tx.Where("block_hash_when_challenge_sent = ? AND message_type = ?", current_block_hash, "storage_challenge_issuance_message").Find(&challenges_issued).Error
		total_challenges_issued = len(challenges_issued)
		var challenges_responded_to []ChallengeMessages
		err = tx.Where("block_hash_when_challenge_sent = ? AND challenge_status = ?", current_block_hash, "Responded").Find(&challenges_responded_to).Error
		total_challenges_responded_to = len(challenges_responded_to)
		var challenges_correct []ChallengeMessages
		err = tx.Where("block_hash_when_challenge_sent = ? AND challenge_status = ?", current_block_hash, "Successful response").Find(&challenges_correct).Error
		total_challenges_correct = len(challenges_correct)
		var challenges_incorrect []ChallengeMessages
		err = tx.Where("block_hash_when_challenge_sent = ? AND challenge_status = ?", current_block_hash, "Failed because of incorrect response").Find(&challenges_incorrect).Error
		total_challenges_incorrect = len(challenges_incorrect)
		var challenges_correct_but_too_slow []ChallengeMessages
		err = tx.Where("block_hash_when_challenge_sent = ? AND challenge_status = ?", current_block_hash, "Failed--Correct but response was too slow").Find(&challenges_correct_but_too_slow).Error
		total_challenges_correct_but_too_slow = len(challenges_correct_but_too_slow)
		var challenges_never_responded_to []ChallengeMessages
		err = tx.Where("block_hash_when_challenge_sent = ? AND challenge_status = ?", current_block_hash, "Failed because response never arrived").Find(&challenges_never_responded_to).Error
		total_challenges_never_responded_to = len(challenges_never_responded_to)
		if total_challenges_issued > 0 {
			challenge_response_success_rate_pct = float32(total_challenges_correct) / float32(total_challenges_issued)
		} else {
			challenge_response_success_rate_pct = float32(1.0)
		}
		var current_block_update PastelBlocks
		current_block_update.Block_hash = current_block_hash
		current_block_update.Total_challenges_issued = uint(total_challenges_issued)
		current_block_update.Total_challenges_responded_to = uint(total_challenges_responded_to)
		current_block_update.Total_challenges_incorrect = uint(total_challenges_incorrect)
		current_block_update.Total_challenges_correct_but_too_slow = uint(total_challenges_correct_but_too_slow)
		current_block_update.Total_challenges_never_responded_to = uint(total_challenges_never_responded_to)
		current_block_update.Challenge_response_success_rate_pct = challenge_response_success_rate_pct
		err = tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "Block_hash"}}, UpdateAll: true}).Create(&current_block_update).Error
	}
}

func RespondToStorageChallenges(responding_masternode_id string, rqsymbol_file_storage_data_folder_path string, current_block_hash string) []string {
	slice_of_message_ids := make([]string, 0)
	db := storage.GetStore().GetDB()
	var pending_challenges []Challenges
	// TODO: handle query error
	db.Where("responding_masternode_id = ? AND challenge_status = ?", responding_masternode_id, "Pending").Find(&pending_challenges)
	slice_of_pending_challenge_ids := make([]string, len(pending_challenges))
	for idx, current_pending_challenge := range pending_challenges {
		slice_of_pending_challenge_ids[idx] = current_pending_challenge.Challenge_id
	}
	var pending_challenge_messages []ChallengeMessages
	// TODO: handle query error
	db.Where("message_type = ? AND challenge_id in ?", "storage_challenge_issuance_message", slice_of_pending_challenge_ids).Find(&pending_challenge_messages)

	for _, current_challenge_message := range pending_challenge_messages {
		x := current_challenge_message
		y := current_challenge_message
		y.Message_type = "storage_challenge_response_message"
		filepath_for_challenge_file_hash := CheckForLocalFilepathForFileHashFunc(responding_masternode_id, x.File_hash_to_challenge, rqsymbol_file_storage_data_folder_path)
		fmt.Println("\nMasternode " + responding_masternode_id + " found a new storage challenge for file hash " + x.File_hash_to_challenge + " (start index: " + fmt.Sprint(x.Challenge_slice_start_index) + "; end index: " + fmt.Sprint(x.Challenge_slice_end_index) + "), responding now!")
		if len(filepath_for_challenge_file_hash) > 0 {
			challenge_file_data, _ := file.ReadFileIntoMemory(filepath_for_challenge_file_hash)
			y.Challenge_response_hash = ComputeHashofFileSlice(challenge_file_data, int(x.Challenge_slice_start_index), int(x.Challenge_slice_end_index))
			challenge_status := "Responded"
			y.Challenge_status = challenge_status
			message_id_input_data := y.Challenging_masternode_id + y.Responding_masternode_id + y.File_hash_to_challenge + y.Challenge_status + y.Message_type + y.Block_hash_when_challenge_sent
			message_id := helper.GetHashFromString(message_id_input_data)
			y.Message_id = message_id
			datetime_challenge_responded_to := time.Now().Format(time.RFC3339)
			y.Datetime_challenge_responded_to = datetime_challenge_responded_to
			UpdateDbWithMessage(y)
			slice_of_message_ids = append(slice_of_message_ids, message_id)
			time_to_respond_to_storage_challenge_in_seconds := helper.ComputeElapsedTimeInSecondsBetweenTwoDatetimes(x.Datetime_challenge_sent, datetime_challenge_responded_to)
			fmt.Println("\nMasternode " + responding_masternode_id + " responded to storage challenge for file hash " + x.File_hash_to_challenge + " in " + fmt.Sprint(time_to_respond_to_storage_challenge_in_seconds) + " seconds!")
		} else {
			fmt.Println("\nMasternode " + responding_masternode_id + " was unable to respond to storage challenge because it did not have the file for file hash " + x.File_hash_to_challenge)
		}
	}
	return slice_of_message_ids
}

func VerifyStorageChallengeResponses(challenging_masternode_id string, rqsymbol_file_storage_data_folder_path string, max_seconds_to_respond_to_storage_challenge int) []string {
	slice_of_message_ids := make([]string, 0)
	db := storage.GetStore().GetDB()
	var responded_challenges []Challenges
	db.Where("challenging_masternode_id = ? AND challenge_status = ?", challenging_masternode_id, "Responded").Find(&responded_challenges)
	slice_of_responded_challenge_ids := make([]string, len(responded_challenges))
	for idx, current_responded_challenge := range responded_challenges {
		slice_of_responded_challenge_ids[idx] = current_responded_challenge.Challenge_id
	}
	var responded_challenge_messages []ChallengeMessages
	db.Where("message_type = ? AND challenge_id in ?", "storage_challenge_response_message", slice_of_responded_challenge_ids).Find(&responded_challenge_messages)
	challenge_status := ""
	for _, current_challenge_response_message := range responded_challenge_messages {
		x := current_challenge_response_message
		y := current_challenge_response_message
		fmt.Println("\nMasternode " + challenging_masternode_id + " found a storage challenge response for file hash " + x.File_hash_to_challenge + " (start index: " + fmt.Sprint(x.Challenge_slice_start_index) + "; end index: " + fmt.Sprint(x.Challenge_slice_end_index) + ") from responding Masternode " + x.Responding_masternode_id + ", verifying now!")
		y.Message_type = "storage_challenge_verification_message"
		filepath_for_challenge_file_hash := CheckForLocalFilepathForFileHashFunc(challenging_masternode_id, x.File_hash_to_challenge, rqsymbol_file_storage_data_folder_path)
		if len(filepath_for_challenge_file_hash) > 0 {
			challenge_file_data, _ := file.ReadFileIntoMemory(filepath_for_challenge_file_hash)
			y.Challenge_slice_correct_hash = ComputeHashofFileSlice(challenge_file_data, int(x.Challenge_slice_start_index), int(x.Challenge_slice_end_index))
			datetime_response_verified := time.Now().Format(time.RFC3339)
			time_to_verify_storage_challenge_in_seconds := helper.ComputeElapsedTimeInSecondsBetweenTwoDatetimes(x.Datetime_challenge_sent, datetime_response_verified)
			if (y.Challenge_slice_correct_hash == y.Challenge_response_hash) && (time_to_verify_storage_challenge_in_seconds <= float64(max_seconds_to_respond_to_storage_challenge)) {
				challenge_status = "Successful response"
				fmt.Println("\nMasternode " + x.Responding_masternode_id + " correctly responded in " + fmt.Sprint(time_to_verify_storage_challenge_in_seconds) + " seconds to a storage challenge for file " + x.File_hash_to_challenge)
			} else if y.Challenge_slice_correct_hash == x.Challenge_response_hash {
				challenge_status = "Failed--Correct but response was too slow"
				fmt.Println("\nMasternode " + x.Responding_masternode_id + " correctly responded in " + fmt.Sprint(time_to_verify_storage_challenge_in_seconds) + " seconds to a storage challenge for file " + x.File_hash_to_challenge + ", but was too slow so failed the challenge anyway!")
			} else {
				challenge_status = "Failed because of incorrect response"
				fmt.Println("\nMasternode " + x.Responding_masternode_id + " failed by incorrectly responding to a storage challenge for file " + x.File_hash_to_challenge)
			}

		} else {
			fmt.Println("\nMasternode " + x.Responding_masternode_id + " was unable to verify the storage challenge response, but it was the fault of the Challenger, so voidng challenge!")
			challenge_status = "Void"
		}
		y.Challenge_status = challenge_status
		y.Datetime_challenge_verified = time.Now().Format(time.RFC3339)
		message_id_input_data := y.Challenging_masternode_id + y.Responding_masternode_id + y.File_hash_to_challenge + y.Challenge_status + y.Message_type + y.Block_hash_when_challenge_sent
		message_id := helper.GetHashFromString(message_id_input_data)
		y.Message_id = message_id
		UpdateDbWithMessage(y)
		slice_of_message_ids = append(slice_of_message_ids, message_id)
	}
	var unresponded_challenges []Challenges
	db.Where("challenging_masternode_id = ? AND challenge_status = ?", challenging_masternode_id, "Pending").Find(&unresponded_challenges)
	slice_of_unresponded_challenge_ids := make([]string, len(unresponded_challenges))
	for idx, current_unresponded_challenge := range unresponded_challenges {
		slice_of_unresponded_challenge_ids[idx] = current_unresponded_challenge.Challenge_id
	}
	var unresponded_challenge_messages []ChallengeMessages
	db.Where("challenge_id in ? AND message_type = ?", slice_of_unresponded_challenge_ids, "storage_challenge_response_message").Find(&unresponded_challenge_messages)
	for _, current_challenge_response_message2 := range unresponded_challenge_messages {
		y := current_challenge_response_message2
		y.Message_type = "storage_challenge_verification_message"
		y.Challenge_status = "Failed because response never arrived"
		y.Datetime_challenge_verified = time.Now().Format(time.RFC3339)
		message_id_input_data := y.Challenging_masternode_id + y.Responding_masternode_id + y.File_hash_to_challenge + y.Challenge_status + y.Message_type + y.Block_hash_when_challenge_sent
		message_id := helper.GetHashFromString(message_id_input_data)
		fmt.Println("\nChallenging Masternode " + challenging_masternode_id + " never got a storage challenge response in the required time from masternode " + y.Responding_masternode_id + " for file hash " + y.File_hash_to_challenge + ", so marking that challenge as a failure!")
		UpdateDbWithMessage(y)
		slice_of_message_ids = append(slice_of_message_ids, message_id)
	}
	return slice_of_message_ids
}

func SimulateDishonestMasternode(dishonest_masternode_id string, approximate_percentage_of_responsible_files_to_ignore float64, rqsymbol_file_storage_data_folder_path string) {
	dishonest_masternode_folder_path := rqsymbol_file_storage_data_folder_path + dishonest_masternode_id + "/"
	slice_of_responsible_filepaths, _ := filepath.Glob(dishonest_masternode_folder_path + "*")
	slice_of_responsible_files_to_randomly_delete := make([]string, 0)
	for _, current_responsible_filepath := range slice_of_responsible_filepaths {
		if rand.Float64() <= approximate_percentage_of_responsible_files_to_ignore {
			slice_of_responsible_files_to_randomly_delete = append(slice_of_responsible_files_to_randomly_delete, current_responsible_filepath)
		}
		for _, current_filepath_to_delete := range slice_of_responsible_files_to_randomly_delete {
			resetReadOnlyFlagAll(current_filepath_to_delete)
			RemoveGlob(current_filepath_to_delete)
		}
	}
}

func MakeDishonestMasternodesDeleteRandomFiles(slice_of_pastel_masternode_ids []string, rqsymbol_file_storage_data_folder_path string) {
	dishonest_masternode_id := "jXlzy0y3L1gYG04DBEZSKI9KV5BReiRzrW5bDBls3M2gtS6R0Ed8MHrEW9hzzgi4aW1taxNzChPSHEgJY4aTbw"
	approximate_percentage_of_responsible_files_to_ignore := 0.35
	fmt.Println("Selected masternode " + dishonest_masternode_id + " to be a dishonest node that deleted " + fmt.Sprint(approximate_percentage_of_responsible_files_to_ignore*100) + " percent of its raptorq symbol files...")
	SimulateDishonestMasternode(dishonest_masternode_id, approximate_percentage_of_responsible_files_to_ignore, rqsymbol_file_storage_data_folder_path)
	dishonest_masternode_id2 := "jXEZVtIEVmSkYw0v8qGjsBrrELBOPuedNYMctelLWSlw6tiVNljFMpZFir30SN9r645tEAKwEAYfKR3o4Ek5YM"
	approximate_percentage_of_responsible_files_to_ignore2 := 0.75
	fmt.Println("Selected masternode " + dishonest_masternode_id2 + " to be a dishonest node that deleted " + fmt.Sprint(approximate_percentage_of_responsible_files_to_ignore2*100) + " percent of its raptorq symbol files...")
	SimulateDishonestMasternode(dishonest_masternode_id2, approximate_percentage_of_responsible_files_to_ignore2, rqsymbol_file_storage_data_folder_path)
	dishonest_masternode_id3 := "jXqBzHsk8P1cuRFrsRkQR5IhPzwFyCxE369KYqFLSITr8l5koLWcabZZDUVltIJ8666bE53G5fbtCz4veU2FCP"
	approximate_percentage_of_responsible_files_to_ignore3 := 0.15
	fmt.Println("Selected masternode " + dishonest_masternode_id3 + " to be a dishonest node that deleted " + fmt.Sprint(approximate_percentage_of_responsible_files_to_ignore3*100) + " percent of its raptorq symbol files...")
	SimulateDishonestMasternode(dishonest_masternode_id3, approximate_percentage_of_responsible_files_to_ignore3, rqsymbol_file_storage_data_folder_path)
	dishonest_masternode_id4 := "jXTwS1eCNDopMUIZAQnvpGlVe9lEnbauoh8TNDRoZcRTJVxCmZu1oSySBM1UwwyHDh7npbn01tZG0q2xyGmVJr"
	approximate_percentage_of_responsible_files_to_ignore4 := 0.05
	fmt.Println("Selected masternode " + dishonest_masternode_id4 + " to be a dishonest node that deleted " + fmt.Sprint(approximate_percentage_of_responsible_files_to_ignore4*100) + " percent of its raptorq symbol files...")
	SimulateDishonestMasternode(dishonest_masternode_id4, approximate_percentage_of_responsible_files_to_ignore4, rqsymbol_file_storage_data_folder_path)
}

func AddNewMasternodeIdsAndFiles(slice_of_new_masternode_ids []string, slice_of_new_file_paths []string, xor_distance_matrix [][]uint64) {
	fmt.Println("Adding " + fmt.Sprint(len(slice_of_new_masternode_ids)) + " new Masternode IDs and " + fmt.Sprint(len(slice_of_new_file_paths)) + " new files to the system...")
	AddXorDistanceMatrixToDb(slice_of_new_masternode_ids, slice_of_new_file_paths, xor_distance_matrix)
	AddFilesToDb(slice_of_new_file_paths)
	AddMasternodesToDb(slice_of_new_masternode_ids)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func AddNIncrementalMasternodeIdsAndKIncrementalFiles(n int, k int, slice_of_new_masternode_ids []string, slice_of_new_file_paths []string, xor_distance_matrix [][]uint64) {
	slice_of_existing_masternode_ids, slice_of_existing_file_hashes := GetCurrentListsOfMasternodeIdsAndFileHashesFromDb()
	slice_of_existing_file_paths := GetSliceOfFilePathsFromSliceOfFileHashes(slice_of_existing_file_hashes)
	incremental_masternode_ids := FindMissingElementsOfAinB(slice_of_new_masternode_ids, slice_of_existing_masternode_ids)
	incremental_file_paths := FindMissingElementsOfAinB(slice_of_new_file_paths, slice_of_existing_file_paths)
	incremental_masternode_count := min(len(incremental_masternode_ids), n)
	incremental_file_path_count := min(len(incremental_file_paths), k)
	AddNewMasternodeIdsAndFiles(incremental_masternode_ids[0:incremental_masternode_count], incremental_file_paths[0:incremental_file_path_count], xor_distance_matrix)
}
