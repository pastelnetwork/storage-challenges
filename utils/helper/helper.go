package helper

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/pastelnetwork/storage-challenges/utils/xordistance"
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

func GetHashFromString(input_string string) string {
	h := sha3.New256()
	h.Write([]byte(input_string))
	sha256_hash_of_input_string := hex.EncodeToString(h.Sum(nil))
	return sha256_hash_of_input_string
}

func GenerateFakeBlockHashes(number_of_blocks_to_make int) []string {
	slice_of_block_hashes := make([]string, number_of_blocks_to_make)
	for ii := range slice_of_block_hashes {
		slice_of_block_hashes[ii] = GetHashFromString(fmt.Sprint(ii))
	}
	return slice_of_block_hashes
}

func ComputeElapsedTimeInSecondsBetweenTwoDatetimes(start_datetime_string string, end_datetime_string string) float64 {
	start_datetime, _ := time.Parse(time.RFC3339, start_datetime_string)
	end_datetime, _ := time.Parse(time.RFC3339, end_datetime_string)
	time_delta := end_datetime.Sub(start_datetime)
	total_seconds_elapsed := time_delta.Seconds()
	return total_seconds_elapsed
}

func ComputeXorDistanceBetweenTwoStrings(string1 string, string2 string) uint64 {
	string1_hash := GetHashFromString(string1)
	string2_hash := GetHashFromString(string2)
	string_1_hash_as_bytes := []byte(string1_hash)
	string_2_hash_as_bytes := []byte(string2_hash)
	xor_distance, _ := xordistance.XORBytes(string_1_hash_as_bytes, string_2_hash_as_bytes)
	xor_distance_as_int := xordistance.BytesToInt(xor_distance)
	xor_distance_as_string := fmt.Sprint(xor_distance_as_int)
	xor_distance_as_string_rescaled := fmt.Sprint(xor_distance_as_string[:len(xor_distance_as_string)-137])
	xor_distance_as_uint64, _ := strconv.ParseUint(xor_distance_as_string_rescaled, 10, 64)
	return xor_distance_as_uint64
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
