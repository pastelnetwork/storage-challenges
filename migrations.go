package main

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
	MasternodeID                     string `gorm:"primaryKey;unique"`
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
	XORDistanceID string `gorm:"primaryKey;unique;column:xor_distance_id"`
	MasternodeID  string
	FileHash      string
	XORDistance   uint64     `gorm:"column:xor_distance"`
	SymbolFile    SymbolFile `gorm:"foreignKey:MasternodeID"`
	Masternode    Masternode `gorm:"foreignKey:FileHash"`
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
	Masternode                     Masternode  `gorm:"foreignKey:ChallengingMasternodeID; foreignKey:RespondingMasternodeID"`
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
	Masternode                    Masternode  `gorm:"foreignKey:ChallengingMasternodeID; foreignKey:RespondingMasternodeID"`
	Challenge                     Challenge   `gorm:"foreignKey:ChallengeID"`
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
			MasternodeID:        "jXlzy0y3L1gYG04DBEZSKI9KV5BReiRzrW5bDBls3M2gtS6R0Ed8MHrEW9hzzgi4aW1taxNzChPSHEgJY4aTbw",
			MasternodeIPAddress: "localhost:9000",
		},
		{
			MasternodeID:        "jXEZVtIEVmSkYw0v8qGjsBrrELBOPuedNYMctelLWSlw6tiVNljFMpZFir30SN9r645tEAKwEAYfKR3o4Ek5YM",
			MasternodeIPAddress: "localhost:9001",
		},
		{
			MasternodeID:        "jXqBzHsk8P1cuRFrsRkQR5IhPzwFyCxE369KYqFLSITr8l5koLWcabZZDUVltIJ8666bE53G5fbtCz4veU2FCP",
			MasternodeIPAddress: "localhost:9002",
		},
		{
			MasternodeID:        "jXTwS1eCNDopMUIZAQnvpGlVe9lEnbauoh8TNDRoZcRTJVxCmZu1oSySBM1UwwyHDh7npbn01tZG0q2xyGmVJr",
			MasternodeIPAddress: "localhost:9003",
		},
	}

	h := sha256.New()
	h.Write([]byte("sample_file_data"))
	filePathHash := h.Sum(nil)

	err = dbTx.Model(&Masternode{}).Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "masternode_id"}}, UpdateAll: true}).Create(mns).Error
	if err != nil {
		return
	}

	xds := []XORDistance{
		{
			XORDistanceID: "xor_distance_id_1",
			FileHash:      base64.StdEncoding.EncodeToString(filePathHash),
			MasternodeID:  "jXlzy0y3L1gYG04DBEZSKI9KV5BReiRzrW5bDBls3M2gtS6R0Ed8MHrEW9hzzgi4aW1taxNzChPSHEgJY4aTbw",
			XORDistance:   1,
		},
		{
			XORDistanceID: "xor_distance_id_2",
			FileHash:      base64.StdEncoding.EncodeToString(filePathHash),
			MasternodeID:  "jXEZVtIEVmSkYw0v8qGjsBrrELBOPuedNYMctelLWSlw6tiVNljFMpZFir30SN9r645tEAKwEAYfKR3o4Ek5YM",
			XORDistance:   2,
		},
		{
			XORDistanceID: "xor_distance_id_3",
			FileHash:      base64.StdEncoding.EncodeToString(filePathHash),
			MasternodeID:  "jXqBzHsk8P1cuRFrsRkQR5IhPzwFyCxE369KYqFLSITr8l5koLWcabZZDUVltIJ8666bE53G5fbtCz4veU2FCP",
			XORDistance:   3,
		},
		{
			XORDistanceID: "xor_distance_id_4",
			FileHash:      base64.StdEncoding.EncodeToString(filePathHash),
			MasternodeID:  "jXTwS1eCNDopMUIZAQnvpGlVe9lEnbauoh8TNDRoZcRTJVxCmZu1oSySBM1UwwyHDh7npbn01tZG0q2xyGmVJr",
			XORDistance:   4,
		},
	}

	err = dbTx.Model(&XORDistance{}).Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "xor_distance_id"}}, UpdateAll: true}).Create(xds).Error

	return err
}
