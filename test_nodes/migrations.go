package testnodes

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommonModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Masternode struct {
	CommonModel
	NodeID                           string `gorm:"primaryKey;unique"`
	MasternodeIPAddress              string `gorm:"column:masternode_ip_address"`
	TotalChallengesIssued            uint
	TotalChallengesRespondedTo       uint
	TotalChallengesCorrect           uint
	TotalChallengesIncorrect         uint
	TotalChallengesCorrectButTooSlow uint
	TotalChallengesNeverRespondedTo  uint
	ChallengeResponseSuccessRatePct  float32
}

type PastelBlock struct {
	CommonModel
	BlockHash                       string `gorm:"primaryKey;unique"`
	BlockNumber                     uint
	TotalChallengesIssued           uint
	TotalChallengesRespondedTo      uint
	TotalChallengesCorrect          uint
	TotalChallengesIncorrect        uint
	TotalChallengeTimeout           uint
	ChallengeResponseSuccessRatePct float32
}

type SymbolFile struct {
	CommonModel
	FileHash               string `gorm:"primaryKey;unique"`
	FileLengthInBytes      uint
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

func AutoMigrate(seeding bool) {
	fmt.Println()
	fmt.Println("*****************************************")
	fmt.Println("*******      START MIGRATION      *******")
	fmt.Println("*****************************************")
	fmt.Println()
	db, err := gorm.Open(sqlite.Open("storage-challenge.sqlite"), &gorm.Config{CreateBatchSize: 1000})
	if err != nil {
		panic("failed to connect database")
	}
	db = db.Debug()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()
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

		if err = dataSeeding(db); err != nil {
			fmt.Println("FAILED TO SEEDING DUMMY DATA:", err)
		}
	}
	fmt.Println()
	fmt.Println("*****************************************")
	fmt.Println("*******   COMPLETED MIGRATIONS    *******")
	fmt.Println("*****************************************")
	fmt.Println()
}

func dataSeeding(db *gorm.DB) (err error) {
	dbTx := db.Begin()
	defer func() {
		if err != nil {
			dbTx.Rollback()
		} else {
			dbTx.Commit()
		}
	}()

	mns := []Masternode{
		{
			NodeID:              "jXlzy0y3L1gYG04DBEZSKI9KV5BReiRzrW5bDBls3M2gtS6R0Ed8MHrEW9hzzgi4aW1taxNzChPSHEgJY4aTbw",
			MasternodeIPAddress: "localhost:9000",
		},
		{
			NodeID:              "jXEZVtIEVmSkYw0v8qGjsBrrELBOPuedNYMctelLWSlw6tiVNljFMpZFir30SN9r645tEAKwEAYfKR3o4Ek5YM",
			MasternodeIPAddress: "localhost:9001",
		},
		{
			NodeID:              "jXqBzHsk8P1cuRFrsRkQR5IhPzwFyCxE369KYqFLSITr8l5koLWcabZZDUVltIJ8666bE53G5fbtCz4veU2FCP",
			MasternodeIPAddress: "localhost:9002",
		},
		{
			NodeID:              "jXTwS1eCNDopMUIZAQnvpGlVe9lEnbauoh8TNDRoZcRTJVxCmZu1oSySBM1UwwyHDh7npbn01tZG0q2xyGmVJr",
			MasternodeIPAddress: "localhost:9003",
		},
	}

	h := sha256.New()
	fileContent := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum"
	h.Write([]byte(fileContent))
	fileHash := base64.StdEncoding.EncodeToString(h.Sum(nil))
	fileLenght := len(fileContent)

	err = dbTx.Model(&Masternode{}).Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "node_id"}}, UpdateAll: true}).Create(mns).Error
	if err != nil {
		return
	}

	xds := []XORDistance{
		{
			XORDistanceID:  "xor_distance_id_1",
			SymbolFileHash: fileHash,

			MasternodeID: "jXlzy0y3L1gYG04DBEZSKI9KV5BReiRzrW5bDBls3M2gtS6R0Ed8MHrEW9hzzgi4aW1taxNzChPSHEgJY4aTbw",
			XORDistance:  1,
		},
		{
			XORDistanceID:  "xor_distance_id_2",
			SymbolFileHash: fileHash,

			MasternodeID: "jXEZVtIEVmSkYw0v8qGjsBrrELBOPuedNYMctelLWSlw6tiVNljFMpZFir30SN9r645tEAKwEAYfKR3o4Ek5YM",
			XORDistance:  2,
		},
		{
			XORDistanceID:  "xor_distance_id_3",
			SymbolFileHash: fileHash,

			MasternodeID: "jXqBzHsk8P1cuRFrsRkQR5IhPzwFyCxE369KYqFLSITr8l5koLWcabZZDUVltIJ8666bE53G5fbtCz4veU2FCP",
			XORDistance:  3,
		},
		{
			XORDistanceID:  "xor_distance_id_4",
			SymbolFileHash: fileHash,

			MasternodeID: "jXTwS1eCNDopMUIZAQnvpGlVe9lEnbauoh8TNDRoZcRTJVxCmZu1oSySBM1UwwyHDh7npbn01tZG0q2xyGmVJr",
			XORDistance:  4,
		},
	}

	err = dbTx.Model(&XORDistance{}).Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "xor_distance_id"}}, UpdateAll: true}).Create(xds).Error
	if err != nil {
		return
	}

	sfs := []SymbolFile{
		{
			FileHash:          fileHash,
			FileLengthInBytes: uint(fileLenght),
			OriginalFilePath:  "test_symbol_files/symbol_file",
		},
	}

	err = dbTx.Model(&SymbolFile{}).Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "file_hash"}}, UpdateAll: true}).Create(sfs).Error
	if err != nil {
		return
	}

	return err
}
