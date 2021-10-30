package model

type ChallengeMessages struct {
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

type PastelBlocks struct {
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

type Masternodes struct {
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

type SymbolFiles struct {
	File_hash                 string
	File_length_in_bytes      uint
	Total_challenges_for_file uint
	Original_file_path        string
}

type XOR_Distance struct {
	Xor_distance_id string
	Masternode_id   string
	File_hash       string
	Xor_distance    uint64
}
