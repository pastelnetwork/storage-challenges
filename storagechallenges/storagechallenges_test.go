package storagechallenges

import "testing"

/* func TestStorageChallenge(t *testing.T) {
	if storagechallenges.TestFunc() != "duku" {
		t.Fatal("wrong answer!")
	}

}
*/

func TestUpdateDbWithMessage(t *testing.T) {
	type args struct {
		storage_challenge_message ChallengeMessages
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"success", args{storage_challenge_message: ChallengeMessages{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UpdateDbWithMessage(tt.args.storage_challenge_message)
		})
	}
}
