package file

import (
	"bufio"
	"encoding/hex"
	"io"
	"os"

	"golang.org/x/crypto/sha3"
)

func ReadFileIntoMemory(input_filepath string) ([]byte, error) {
	file, err := os.Open(input_filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil, statsErr
	}
	var size int64 = stats.Size()
	bytes := make([]byte, size)
	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)
	return bytes, err
}

func GetHashAndSizeFromFilePath(path_to_file string) (hash string, size uint64, err error) {
	f, err := os.Open(path_to_file)
	if err != nil {
		return "", 0, err
	}
	defer f.Close()
	fileInfo, err := f.Stat()
	if err != nil {
		return "", 0, err
	}
	h := sha3.New256()

	if _, err := io.Copy(h, f); err != nil {
		return "", 0, err
	}
	return hex.EncodeToString(h.Sum(nil)), uint64(fileInfo.Size()), nil
}
