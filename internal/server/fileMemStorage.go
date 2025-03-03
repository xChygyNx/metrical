package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/xChygyNx/metrical/internal/server/types"
)

const (
	filePem = 0o600
)

func fileDump(fileName string, period time.Duration, storage *types.MemStorage) (err error) {
	if err != nil {
		return
	}
	ticker := time.NewTicker(period)
	defer ticker.Stop()
	for range ticker.C {
		err := writeMetricStorageFile(fileName, storage)
		if err != nil {
			return fmt.Errorf("error in write data in metric storage file: %w", err)
		}
	}
	return
}

func writeMetricStorageFile(absStorageFilePath string, storage *types.MemStorage) (err error) {
	file, err := os.OpenFile(absStorageFilePath, os.O_WRONLY|os.O_CREATE, filePem)
	if err != nil {
		return fmt.Errorf("error in open file %s: %w", absStorageFilePath, err)
	}
	defer func() {
		err = file.Close()
	}()
	data, err := json.Marshal(*storage)
	if err != nil {
		return fmt.Errorf("error in marshal data for record in fille: %w", err)
	}
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("error write data in fille %s: %w", absStorageFilePath, err)
	}
	return
}

func retryFileWrite(absStorageFilePath string, storage *types.MemStorage, retryCount int) (err error) {
	delays := make([]time.Duration, 0, retryCount)
	delays = append(delays, 0*time.Second)
	for i := 1; i < retryCount-1; i++ {
		delays = append(delays, time.Duration(2*i-1)*time.Second)
	}

	for i := 0; i < retryCount; i++ {
		time.Sleep(delays[i])
		err = writeMetricStorageFile(absStorageFilePath, storage)
		if err == nil {
			break
		}
	}
	return
}

func restoreMetricStore(fileName string, storage *types.MemStorage) (err error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, filePem)
	if err != nil {
		return
	}
	defer func() {
		err = file.Close()
	}()

	sc := bufio.NewScanner(file)
	if !sc.Scan() {
		return fmt.Errorf("error scan metric storage file: %w", sc.Err())
	}
	data := sc.Bytes()

	err = json.Unmarshal(data, storage)
	if err != nil {
		return
	}

	return
}
