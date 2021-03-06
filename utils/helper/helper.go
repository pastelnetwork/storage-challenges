package helper

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/sha3"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func FillSliceWithFunctionOutput(vs []string, f func() string) []string {
	vsm := make([]string, len(vs))
	for i := range vs {
		vsm[i] = f()
	}
	return vsm
}

func GetHashFromString(inputString string) string {
	h := sha3.New256()
	h.Write([]byte(inputString))
	sha256HashOfInputString := hex.EncodeToString(h.Sum(nil))
	return sha256HashOfInputString
}

func GenerateFakeBlockHashes(numberOfBlockToMake int) []string {
	sliceOfBlockHash := make([]string, numberOfBlockToMake)
	for ii := range sliceOfBlockHash {
		sliceOfBlockHash[ii] = GetHashFromString(fmt.Sprint(ii))
	}
	return sliceOfBlockHash
}

func ComputeElapsedTimeInSecondsBetweenTwoDatetimes(start, end int64) float64 {
	return float64(end - start)
}

func GenerateFakePastelMnID() string {
	fake_id_prefix := "jX"
	fake_id := fake_id_prefix + GenerateRandomAlphaNumericString(84)
	return fake_id
}

func ApplyFunctionToElementsInSlice(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func GenerateRandomAlphaNumericString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GetMasternodeIds(number_of_masternode_ids_to_make int) []string {
	slice_of_masternode_ids := make([]string, number_of_masternode_ids_to_make)
	slice_of_masternode_ids = FillSliceWithFunctionOutput(slice_of_masternode_ids, GenerateFakePastelMnID)
	return slice_of_masternode_ids
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

func SliceContainsString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
