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

//func getMemStorageFileAbsPath(fileName string) (string, error) {
//	dirPath, err := filepath.Abs("./memory_metrics")
//	if err != nil {
//		return "", fmt.Errorf("can't get absolute path: %w", err)
//	}
//
//	if _, err = os.Stat(dirPath); errors.Is(err, os.ErrNotExist) {
//		err = os.Mkdir(dirPath, dirPerm)
//		if err != nil {
//			return "", fmt.Errorf("can't create directory %s: %w", dirPath, err)
//		}
//	} else if err != nil {
//		return "", fmt.Errorf("error in search directory %s: %w", dirPath, err)
//	}
//	return filepath.Join(dirPath, fileName), nil
//}

func fileDump(fileName string, period time.Duration, storage *types.MemStorage) (err error) {
	//storageFilePath, err := getMemStorageFileAbsPath(fileName)
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

func restoreMetricStore(fileName string, storage *types.MemStorage) (err error) {
	//storageFilePath, err := getMemStorageFileAbsPath(fileName)
	//if err != nil {
	//	return
	//}

	if _, err = os.Stat(fileName); err != nil {
		return fmt.Errorf("can't find metrics storage file: %w", err)
	}

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
