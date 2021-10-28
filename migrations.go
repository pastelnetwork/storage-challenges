package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func AutoMigrate() {
	db, err := gorm.Open(sqlite.Open("go_pastel_storage_challenges.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	type Masternodes struct {
		gorm.Model
		Masternode_id                         string `gorm:"primaryKey;unique"`
		Masternode_ip_address                 string
		Total_challenges_issued               uint
		Total_challenges_responded_to         uint
		Total_challenges_correct              uint
		Total_challenges_incorrect            uint
		Total_challenges_correct_but_too_slow uint
		Total_challenges_never_responded_to   uint
		Challenge_response_success_rate_pct   float32
	}

	type PastelBlocks struct {
		gorm.Model
		Block_hash                            string `gorm:"primaryKey;unique"`
		Block_number                          uint
		Total_challenges_issued               uint
		Total_challenges_responded_to         uint
		Total_challenges_correct              uint
		Total_challenges_incorrect            uint
		Total_challenges_correct_but_too_slow uint
		Total_challenges_never_responded_to   uint
		Challenge_response_success_rate_pct   float32
	}

	type SymbolFiles struct {
		gorm.Model
		File_hash                 string `gorm:"primaryKey;unique"`
		File_length_in_bytes      uint
		Total_challenges_for_file uint
		Original_file_path        string
	}

	type XOR_Distance struct {
		gorm.Model
		Xor_distance_id string `gorm:"primaryKey;unique"`
		Masternode_id   string
		File_hash       string
		Xor_distance    uint64
		SymbolFiles     SymbolFiles `gorm:"foreignKey:Masternode_id"`
		Masternodes     Masternodes `gorm:"foreignKey:File_hash"`
	}

	type Challenges struct {
		gorm.Model
		Challenge_id                       string `gorm:"primaryKey;unique"`
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
		PastelBlocks                       PastelBlocks `gorm:"foreignKey:Block_hash_when_challenge_sent"`
		SymbolFiles                        SymbolFiles  `gorm:"foreignKey:File_hash_to_challenge"`
		Masternodes                        Masternodes  `gorm:"foreignKey:Challenging_masternode_id; foreignKey:Responding_masternode_id"`
	}

	type ChallengeMessages struct {
		gorm.Model
		Message_id                      string `gorm:"primaryKey;unique"`
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
		PastelBlocks                    PastelBlocks `gorm:"foreignKey:Block_hash_when_challenge_sent"`
		SymbolFiles                     SymbolFiles  `gorm:"foreignKey:File_hash_to_challenge"`
		Masternodes                     Masternodes  `gorm:"foreignKey:Challenging_masternode_id; foreignKey:Responding_masternode_id"`
		Challenges                      Challenges   `gorm:"foreignKey:Challenge_id"`
	}

	db.AutoMigrate(
		&Masternodes{},
		&PastelBlocks{},
		&SymbolFiles{},
		&XOR_Distance{},
		&Challenges{},
		&ChallengeMessages{},
	)
}
