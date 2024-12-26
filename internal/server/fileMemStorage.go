package server

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/xChygyNx/metrical/internal/server/types"
)

func getMemStorageFileAbsPath(fileName string) (string, error) {
	dirPath, err := filepath.Abs("./memory_metrics")
	if err != nil {
		return "", err
	}

	if _, err = os.Stat(dirPath); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(dirPath, 0o777)
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}
	return filepath.Join(dirPath, fileName), nil
}

func fileDump(fileName string, period time.Duration, storage *types.MemStorage) (err error) {
	storageFilePath, err := getMemStorageFileAbsPath(fileName)
	if err != nil {
		return
	}
	ticker := time.NewTicker(period)
	defer ticker.Stop()
	select {
	case <-ticker.C:
		err := writeMetricStorageFile(storageFilePath, storage)
		if err != nil {
			return fmt.Errorf("error in write data in metric storage file: %w", err)
		}
	}
	return
}

func writeMetricStorageFile(absStorageFilePath string, storage *types.MemStorage) (err error) {
	file, err := os.OpenFile(absStorageFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func() {
		err = file.Close()
	}()
	data, err := json.Marshal(*storage)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return
}

func restoreMetricStore(fileName string, storage *types.MemStorage) (err error) {
	storageFilePath, err := getMemStorageFileAbsPath(fileName)
	if err != nil {
		return
	}

	if _, err = os.Stat(storageFilePath); err != nil {
		return fmt.Errorf("can't find metrics storage file: %w", err)
	}

	file, err := os.OpenFile(storageFilePath, os.O_RDONLY, 0o666)
	if err != nil {
		return
	}
	defer func() {
		err = file.Close()
	}()

	sc := bufio.NewScanner(file)
	if !sc.Scan() {
		return sc.Err()
	}
	data := sc.Bytes()

	err = json.Unmarshal(data, storage)
	if err != nil {
		return
	}

	return
}
