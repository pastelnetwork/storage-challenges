package testnodes

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/pastelnetwork/storage-challenges/utils/file"
	"github.com/pastelnetwork/storage-challenges/utils/helper"
	"github.com/pastelnetwork/storage-challenges/utils/xordistance"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommonModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Masternode struct {
	CommonModel
	NodeID                          string `gorm:"primaryKey;unique"`
	MasternodeIPAddress             string `gorm:"column:masternode_ip_addre00"`
	TotalChallengesIssued           uint
	TotalChallengesRespondedTo      uint
	TotalChallengesCorrect          uint
	TotalChallengesIncorrect        uint
	TotalChallengesTimeout          uint
	ChallengeResponseSuccessRatePct float32
}

type PastelBlock struct {
	CommonModel
	BlockHash                       string `gorm:"primaryKey;unique"`
	BlockNumber                     uint
	TotalChallengesIssued           uint
	TotalChallengesRespondedTo      uint
	TotalChallengesCorrect          uint
	TotalChallengesIncorrect        uint
	TotalChallengesTimeout          uint
	ChallengeResponseSuccessRatePct float32
}

type SymbolFile struct {
	CommonModel
	FileHash               string `gorm:"primaryKey;unique"`
	FileLengthInBytes      uint64
	TotalChallengesForFile uint
	OriginalFilePath       string
}

type XORDistance struct {
	CommonModel
	XORDistanceID  string `gorm:"primaryKey;unique;column:xor_distance_id"`
	MasternodeID   string
	SymbolFileHash string
	XORDistance    uint64     `gorm:"column:xor_distance"`
	SymbolFile     SymbolFile `gorm:"foreignKey:SymbolFileHash;references:FileHash"`
	Masternode     Masternode `gorm:"foreignKey:MasternodeID;references:NodeID"`
}

func (XORDistance) TableName() string { return "xor_distances" }

type Challenge struct {
	CommonModel
	ChallengeID                    string `gorm:"primaryKey;unique"`
	ChallengeStatus                string
	TimestampChallengeSent         int64
	TimestampChallengeRespondedTo  int64
	TimestampChallengeVerified     int64
	BlockHashWhenChallengeSent     string
	ChallengeResponseTimeInSeconds float64
	ChallengingMasternodeID        string
	RespondingMasternodeID         string
	FileHashToChallenge            string
	ChallengeSliceStartIndex       uint64
	ChallengeSliceEndIndex         uint64
	ChallengeSliceCorrectHash      string
	ChallengeResponseHash          string
	PastelBlock                    PastelBlock `gorm:"foreignKey:BlockHashWhenChallengeSent"`
	SymbolFile                     SymbolFile  `gorm:"foreignKey:FileHashToChallenge"`
	ChallengingMasternode          Masternode  `gorm:"foreignKey:ChallengingMasternodeID"`
	RespondingMasternode           Masternode  `gorm:"foreignKey:RespondingMasternodeID"`
}

type ChallengeMessage struct {
	CommonModel
	MessageID                     string `gorm:"primaryKey;unique"`
	MessageType                   string
	ChallengeStatus               string
	TimestampChallengeSent        int64
	TimestampChallengeRespondedTo int64
	TimestampChallengeVerified    int64
	BlockHashWhenChallengeSent    string
	ChallengingMasternodeID       string
	RespondingMasternodeID        string
	FileHashToChallenge           string
	ChallengeSliceStartIndex      uint64
	ChallengeSliceEndIndex        uint64
	ChallengeSliceCorrectHash     string
	ChallengeResponseHash         string
	ChallengeID                   string
	PastelBlock                   PastelBlock `gorm:"foreignKey:BlockHashWhenChallengeSent"`
	SymbolFile                    SymbolFile  `gorm:"foreignKey:FileHashToChallenge"`
	ChallengingMasternode         Masternode  `gorm:"foreignKey:ChallengingMasternodeID"`
	RespondingMasternode          Masternode  `gorm:"foreignKey:RespondingMasternodeID"`
	Challenge                     Challenge   `gorm:"association_foreignKey:ChallengeID;references:ChallengeID"`
}

func AutoMigrate(db *gorm.DB, seeding bool) {
	fmt.Println()
	fmt.Println("*****************************************")
	fmt.Println("*******      START MIGRATION      *******")
	fmt.Println("*****************************************")
	fmt.Println()
	db.AutoMigrate(
		&Masternode{},
		&PastelBlock{},
		&SymbolFile{},
		&XORDistance{},
		&Challenge{},
		&ChallengeMessage{},
	)

	if seeding {
		fmt.Println()
		fmt.Println("*****************************************")
		fmt.Println("******* START SEEDING DUMMY DATA  *******")
		fmt.Println("*****************************************")
		fmt.Println()

		if err := dataSeeding(db); err != nil {
			fmt.Println("FAILED TO SEEDING DUMMY DATA:", err)
		}
	}
	fmt.Println()
	fmt.Println("*****************************************")
	fmt.Println("*******   COMPLETED MIGRATIONS    *******")
	fmt.Println("*****************************************")
	fmt.Println()
}

type masternodes []Masternode

func (ms masternodes) ListIDs() []string {
	var ids = make([]string, len(ms))
	for idx, m := range ms {
		ids[idx] = m.NodeID
	}
	return ids
}

type symbolFiles []SymbolFile

func (ms symbolFiles) ListFileHashes() []string {
	var ids = make([]string, len(ms))
	for idx, m := range ms {
		ids[idx] = m.FileHash
	}
	return ids
}

func GetListMasternodeIDsFromMasternodes(mns []Masternode) []string {
	return masternodes(mns).ListIDs()
}

var (
	// first 2 nodes
	mns = []Masternode{
		{
			NodeID:              "jXlzy0y3L1gYG04DBEZSKI9KV5BReiRzrW5bDBls3M2gtS6R0Ed8MHrEW9hzzgi4aW1taxNzChPSHEgJY4aTbw",
			MasternodeIPAddress: "node0:9000",
		},
		{
			NodeID:              "jXEZVtIEVmSkYw0v8qGjsBrrELBOPuedNYMctelLWSlw6tiVNljFMpZFir30SN9r645tEAKwEAYfKR3o4Ek5YM",
			MasternodeIPAddress: "node1:9000",
		},
	}

	// 4 dishonest nodes
	dishonestMns = masternodes{
		{
			NodeID:              "jX7RRUiOCNmoggpO67DOAH5An9raJspnY2noBe3UaAlCMqOEo2QQukhI8w0jjiAA78xpwlFc8ucpcV77pjw9Jm",
			MasternodeIPAddress: "node2:9000",
		},
		{
			NodeID:              "jXoIquQRCdRrnjOClioRrSdG6pGyqG3audIQrVwIc6OgR3FFa90WemZ1xuylKjUBMj3gZpL69GT2fdJV99jB81",
			MasternodeIPAddress: "node3:9000",
		},
		{
			NodeID:              "jXAXIVujFd2urNsR3mF1YogDlSKaJVdNx2bXWEo3tZukaICMYKFMBoJUcLeWIHyA1NWXHU9rCp1I32OxY6bKcr",
			MasternodeIPAddress: "node4:9000",
		},
		{
			NodeID:              "jXqsiabBVA07RRwaLfhKu4sQ4SCKSgp7TIcUufwDVZvBTdAD2mihLfdG0H7ZhHQTK2LAbKBGMGwlDPInKWsBMy",
			MasternodeIPAddress: "node5:9000",
		},
	}

	newMasternodes = masternodes{
		{
			NodeID:              "jXqBzHsk8P1cuRFrsRkQR5IhPzwFyCxE369KYqFLSITr8l5koLWcabZZDUVltIJ8666bE53G5fbtCz4veU2FCP",
			MasternodeIPAddress: "node6:9000",
		},
		{
			NodeID:              "jXTwS1eCNDopMUIZAQnvpGlVe9lEnbauoh8TNDRoZcRTJVxCmZu1oSySBM1UwwyHDh7npbn01tZG0q2xyGmVJr",
			MasternodeIPAddress: "node7:9000",
		},
		{
			NodeID:              "jXyCj6J8UXeughB7olBCOBtRylx8fuEESzMcsIgdWGkMbx89J9bY1FaYtMbftCTev9206SI0jY5zIVyELvcoGh",
			MasternodeIPAddress: "node8:9000",
		},
		{
			NodeID:              "jXyFFTa8UAGvMRRpoZWa6L0s4dGVVIAyKEobPCeagrljgshH5eGQTX5nh0z3azAgLlVIoj6aznno6Vq0tiFkfQ",
			MasternodeIPAddress: "node9:9000",
		},
		{
			NodeID:              "jXN0gNcapBcqrMYj28s3QS42txVNEHLvizx48FqRQusivXDtRPqiwXRk3zJ2rHQj0CXa1arrp8eWLCO84n5RIL",
			MasternodeIPAddress: "node10:9000",
		},
		{
			NodeID:              "jXderFvKIhkQyaLV134WNDkV9B5lSRqthT6aU35prg8z3snszlW9bh2A5S78c7oiI9ROZKGb9TbFHzvyuF4X3V",
			MasternodeIPAddress: "node11:9000",
		},
	}
	mapApproximatePercentageOfDishonestMasternodeToResponsibleFilesToIgnore = map[string]int{
		"jX7RRUiOCNmoggpO67DOAH5An9raJspnY2noBe3UaAlCMqOEo2QQukhI8w0jjiAA78xpwlFc8ucpcV77pjw9Jm": 41,
		"jXoIquQRCdRrnjOClioRrSdG6pGyqG3audIQrVwIc6OgR3FFa90WemZ1xuylKjUBMj3gZpL69GT2fdJV99jB81": 87,
		"jXAXIVujFd2urNsR3mF1YogDlSKaJVdNx2bXWEo3tZukaICMYKFMBoJUcLeWIHyA1NWXHU9rCp1I32OxY6bKcr": 47,
		"jXqsiabBVA07RRwaLfhKu4sQ4SCKSgp7TIcUufwDVZvBTdAD2mihLfdG0H7ZhHQTK2LAbKBGMGwlDPInKWsBMy": 29,
	}
)

func init() {
	for masternodeID, ignorePct := range mapApproximatePercentageOfDishonestMasternodeToResponsibleFilesToIgnore {
		log.Println(masternodeID, "---", ignorePct)
	}
}

func dataSeeding(db *gorm.DB) (err error) {
	var wg sync.WaitGroup
	wg.Add(1)
	go removeRandomSymbolFiles(&wg)
	var symbolFilesPath []string
	var symbolFilesFolderPath = "sample_raptorq_symbol_files"
	symbolFilesPath, err = filepath.Glob(symbolFilesFolderPath + "/*")
	if err != nil {
		log.Panicln("filepath.Glob", err)
		return
	}
	log.Printf("found %d symbol files in path %s", len(symbolFilesPath), symbolFilesFolderPath)

	allMns := append(mns, dishonestMns...)
	tx := db.Begin()
	err = tx.Model(&Masternode{}).Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "node_id"}}, UpdateAll: true}).Create(allMns).Error
	if err != nil {
		log.Printf("Inserting %d masternodes failed, doing rollback...", len(mns))
		tx.Rollback()
		return
	}
	tx.Commit()
	log.Printf("Inserted %d masternodes", len(allMns))

	log.Println("NUMBER OF MASTERNODES ", len(mns), "NUMBER OF DISHONEST MASTERNODE ", len(dishonestMns))

	var maxProcessingSymbolFilesPerConcurent = 100
	for cnt := 0; cnt < len(symbolFilesPath); cnt += maxProcessingSymbolFilesPerConcurent {
		wg.Add(1)
		if cnt+maxProcessingSymbolFilesPerConcurent < len(symbolFilesPath) {
			go insertSymbolFilesAndXORDistanceToMasternodes(symbolFilesPath[cnt:cnt+maxProcessingSymbolFilesPerConcurent], mns, dishonestMns, db, &wg)
		} else {
			go insertSymbolFilesAndXORDistanceToMasternodes(symbolFilesPath[cnt:len(symbolFilesPath)-1], mns, dishonestMns, db, &wg)
		}
	}

	wg.Wait()

	return err
}

func insertSymbolFilesAndXORDistanceToMasternodes(symbolFilesPath []string, masternodes, dishonestMasternodes []Masternode, db *gorm.DB, wg *sync.WaitGroup) {
	defer wg.Done()
	var symbolFiles = make([]SymbolFile, 0)
	var mapMasternodeToRelatedXORDistances = make(map[string][]XORDistance)
	var totalXORDistancesInserts = 0
	for _, masternode := range append(masternodes, dishonestMasternodes...) {
		mapMasternodeToRelatedXORDistances[masternode.NodeID] = make([]XORDistance, 0)
	}
	for _, filePath := range symbolFilesPath {
		fileHash, size, err := file.GetHashAndSizeFromFilePath(filePath)
		if err != nil {
			log.Printf("ignoring file '%s' because cannot generate file hash", filePath)
			continue
		}
		symbolFiles = append(symbolFiles, SymbolFile{OriginalFilePath: filePath, FileHash: fileHash, FileLengthInBytes: size})
		for _, masternode := range masternodes {
			mapMasternodeToRelatedXORDistances[masternode.NodeID] = append(mapMasternodeToRelatedXORDistances[masternode.NodeID], XORDistance{
				XORDistanceID:  helper.GetHashFromString(masternode.NodeID + fileHash),
				MasternodeID:   masternode.NodeID,
				SymbolFileHash: fileHash,
				XORDistance:    xordistance.ComputeXorDistanceBetweenTwoStrings(fileHash, masternode.NodeID),
			})
		}
		for _, dishonestMasternode := range dishonestMasternodes {
			// dishonest masternode containing around x% of total symbol files
			mapMasternodeToRelatedXORDistances[dishonestMasternode.NodeID] = append(mapMasternodeToRelatedXORDistances[dishonestMasternode.NodeID], XORDistance{
				XORDistanceID:  helper.GetHashFromString(dishonestMasternode.NodeID + fileHash),
				MasternodeID:   dishonestMasternode.NodeID,
				SymbolFileHash: fileHash,
				XORDistance:    xordistance.ComputeXorDistanceBetweenTwoStrings(fileHash, dishonestMasternode.NodeID),
			})
		}
	}
	tx := db.Begin()
	err := tx.Model(&SymbolFile{}).Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "file_hash"}}, UpdateAll: true}).Create(symbolFiles).Error
	if err != nil {
		log.Printf("Inserting %d symbol files failed, doing rollback...", len(symbolFiles))
		tx.Rollback()
		return
	}
	log.Printf("Inserting %d symbol files", len(symbolFiles))

	for nodeID, xorDistances := range mapMasternodeToRelatedXORDistances {
		err := tx.Model(&XORDistance{}).Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "xor_distance_id"}}, UpdateAll: true}).Create(xorDistances).Error
		if err != nil {
			log.Printf("Inserting %d xor distances to masternode %s failed, doing rollback...", len(xorDistances), nodeID)
			tx.Rollback()
			return
		}
		log.Printf("Inserting %d xor distances of masternode %s", len(xorDistances), nodeID)
		totalXORDistancesInserts += len(xorDistances)
	}

	log.Printf("Inserted %d symbol files with %d xor distances related to %d masternodes", len(symbolFiles), totalXORDistancesInserts, len(masternodes)+len(dishonestMasternodes))
	tx.Commit()
}

func RemoveMasternodesAndSymbolFiles(db *gorm.DB) (err error) {
	var listExistingMasternodes = masternodes{}
	if err = db.Model(&Masternode{}).Find(&listExistingMasternodes).Error; err != nil {
		log.Printf("Cannot query list existing masternode from database: %v", err)
		return err
	}
	var listExistingSymbolFiles = symbolFiles{}
	if err = db.Model(&SymbolFile{}).Find(&listExistingSymbolFiles).Error; err != nil {
		log.Printf("Cannot query list symbol files from database: %v", err)
		return err
	}

	removingMasternodes := listExistingMasternodes[:len(listExistingMasternodes)/3]
	removingSymbolFiles := listExistingSymbolFiles[:len(listExistingSymbolFiles)/3]

	tx := db.Begin()
	if err = tx.Delete(&XORDistance{}, "masternode_id IN ?", removingMasternodes.ListIDs()).Error; err != nil {
		tx.Rollback()
		log.Printf("Cannot delete xor-distance data related to %d nodes from database: %v", len(removingMasternodes), err)
		return err
	}
	log.Printf("Deleteing xor-distance data related to %d nodes from database", len(removingMasternodes))
	if err = tx.Delete(&Masternode{}, "node_id IN ?", removingMasternodes.ListIDs()).Error; err != nil {
		tx.Rollback()
		log.Printf("Cannot delete %d masternodes from database: %v", len(removingMasternodes), err)
		return err
	}
	log.Printf("Deleteing %d masternodes from database", len(removingMasternodes))
	tx.Commit()

	tx = db.Begin()
	if err = tx.Delete(&XORDistance{}, "symbol_file_hash IN ?", removingSymbolFiles.ListFileHashes()).Error; err != nil {
		tx.Rollback()
		log.Printf("Cannot delete xor-distance data related to %d symbol files from database: %v", len(removingSymbolFiles), err)
		return err
	}
	log.Printf("Deleteing xor-distance data related to %d symbol files from database", len(removingSymbolFiles))
	if err = tx.Delete(&SymbolFile{}, "file_hash IN ?", removingSymbolFiles.ListFileHashes()).Error; err != nil {
		tx.Rollback()
		log.Printf("Cannot delete %d symbol files from database: %v", len(removingSymbolFiles), err)
		return err
	}
	log.Printf("Deleteing %d symbol files from database", len(removingSymbolFiles))
	tx.Commit()
	return nil
}

func AddNIncrementalMasternodesAndKIncrementalSymbolFiles(n, k int, db *gorm.DB) (err error) {
	var newSymbolFilesPaths []string
	var newSymbolFilesFolderPath = "incremental_raptorq_symbol_files"

	newSymbolFilesPaths, err = filepath.Glob(newSymbolFilesFolderPath + "/*")
	if err != nil {
		log.Panicln("filepath.Glob", err)
		return
	}
	log.Printf("found %d symbol files in path %s", len(newSymbolFilesPaths), newSymbolFilesFolderPath)

	var listExistingFilePaths []string
	if err = db.Model(&SymbolFile{}).Distinct("original_file_path").Find(&listExistingFilePaths).Error; err != nil {
		log.Printf("Cannot query list existing file hash from database: %v", err)
		return err
	}
	log.Printf("Found %d existing original symbol file path from database", len(listExistingFilePaths))

	var listExistingMasternodes = masternodes{}
	if err = db.Model(&Masternode{}).Find(&listExistingMasternodes).Error; err != nil {
		log.Printf("Cannot query list existing masternode from database: %v", err)
		return err
	}
	log.Printf("Found %d existing masternodes from database", len(listExistingMasternodes))

	incrementalSymbolFilePaths := helper.FindMissingElementsOfAinB(listExistingFilePaths, newSymbolFilesPaths)

	var mapMasternodes = make(map[string]Masternode)
	for _, masternode := range append(listExistingMasternodes, newMasternodes...) {
		mapMasternodes[masternode.NodeID] = masternode
	}
	incrementalMasternodeIDs := helper.FindMissingElementsOfAinB(listExistingMasternodes.ListIDs(), newMasternodes.ListIDs())

	incrementalMasternodeCount := min(len(incrementalMasternodeIDs), n)
	incrementalSymbolFilePathCount := min(len(incrementalSymbolFilePaths), k)

	var incrementalMasternodes = []Masternode{}
	for _, incrementalMasternodeID := range incrementalMasternodeIDs[:incrementalMasternodeCount] {
		incrementalMasternodes = append(incrementalMasternodes, mapMasternodes[incrementalMasternodeID])
	}

	tx := db.Begin()
	err = tx.Model(&Masternode{}).Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "node_id"}}, UpdateAll: true}).Create(incrementalMasternodes).Error
	if err != nil {
		log.Printf("Inserting %d masternodes failed, doing rollback...", len(mns))
		tx.Rollback()
		return
	}
	tx.Commit()
	log.Printf("Inserted %d masternodes", len(incrementalMasternodes))

	var wg = sync.WaitGroup{}
	wg.Add(2)
	// inserts all file hash related with only incremental masternodes
	go insertSymbolFilesAndXORDistanceToMasternodes(append(listExistingFilePaths, incrementalSymbolFilePaths[:incrementalSymbolFilePathCount]...), incrementalMasternodes, []Masternode{}, db, &wg)

	onlyExistingHonestMasternodeIDs := helper.FindMissingElementsOfAinB(listExistingMasternodes.ListIDs(), dishonestMns.ListIDs())
	var onlyExistingHonestMasternodes = masternodes{}
	for _, incrementalMasternodeID := range onlyExistingHonestMasternodeIDs {
		onlyExistingHonestMasternodes = append(onlyExistingHonestMasternodes, mapMasternodes[incrementalMasternodeID])
	}
	// inserts all honest masternode and dishonest masternode with only incremental file paths
	go insertSymbolFilesAndXORDistanceToMasternodes(incrementalSymbolFilePaths[:incrementalSymbolFilePathCount], onlyExistingHonestMasternodes, dishonestMns, db, &wg)

	wg.Wait()

	return err
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func AddPastelBlock(blockIdx int32, db *gorm.DB) (err error) {
	tx := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	var pastelBlock = PastelBlock{
		BlockHash:   helper.GetHashFromString("mock block hash " + fmt.Sprint(blockIdx)),
		BlockNumber: uint(blockIdx),
	}

	return tx.Model(&PastelBlock{}).Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "block_hash"}}, DoNothing: true}).Create(&pastelBlock).Error
}

func GetLastBlockNumer(db *gorm.DB) int32 {
	var blocknum int32
	db.Model(&PastelBlock{}).Select("block_number").Order("block_number DESC").Limit(1).Find(&blocknum)
	return blocknum
}

func GetPastelBlockHash(blockIdx int32) string {
	return helper.GetHashFromString("mock block hash " + fmt.Sprint(blockIdx))
}

func GetMasternodes(db *gorm.DB) []Masternode {
	var masternodes = []Masternode{}
	db.Model(&Masternode{}).Find(&masternodes)
	return masternodes
}

func removeRandomSymbolFiles(wg *sync.WaitGroup) {
	defer wg.Done()
	var symbolFolderPath = "test_nodes/sample_raptorq_symbol_files"
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	var folderPathSuffixes = []int{29, 41, 47, 87}
	for _, folderPathIgnoreSf := range folderPathSuffixes {
		symbolFilesPaths, err := filepath.Glob(symbolFolderPath + fmt.Sprintf("_%d", folderPathIgnoreSf) + "/*")
		if err != nil {
			log.Panicln("filepath.Glob", err)
			return
		}
		log.Printf("found %d symbol files in path %s", len(symbolFilesPaths), symbolFolderPath)

		for _, filePath := range symbolFilesPaths {
			rate := r.Intn(100)
			if rate <= folderPathIgnoreSf {
				err = os.Remove(filePath)
				if err != nil {
					log.Printf("could not remove file %s", filePath)
				}
			}
		}
	}

	var incrementalSymbolFolderPath = "test_nodes/incremental_raptorq_symbol_files"
	for _, folderPathIgnoreSf := range folderPathSuffixes {
		incrementalSymbolFilesPaths, err := filepath.Glob(incrementalSymbolFolderPath + fmt.Sprintf("_%d", folderPathIgnoreSf) + "/*")
		if err != nil {
			log.Panicln("filepath.Glob", err)
			return
		}
		log.Printf("found %d incremental symbol files in path %s", len(incrementalSymbolFilesPaths), incrementalSymbolFolderPath)

		for _, filePath := range incrementalSymbolFilesPaths {
			rate := r.Intn(100)
			if rate <= folderPathIgnoreSf {
				err = os.Remove(filePath)
				if err != nil {
					log.Printf("could not remove file %s", filePath)
				}
			}
		}
	}
}
