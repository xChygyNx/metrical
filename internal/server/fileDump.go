package server

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/xChygyNx/metrical/internal/server/types"
)

func fileDump(fileName string, period time.Duration, storage *types.MemStorage) (err error) {
	dirPath, err := filepath.Abs("./memory_metrics")
	if err != nil {
		return
	}

	if _, err = os.Stat(dirPath); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(dirPath, 0o777)
		if err != nil {
			return err
		}
	} else if err != nil {
		return
	}
	filePath := filepath.Join(dirPath, fileName)
	ticker := time.NewTicker(period)
	defer ticker.Stop()
	select {
	case <-ticker.C:
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
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
	}
	return
}
