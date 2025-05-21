// Package server отвечает за прием метрик от агента и их хранения в соответствии
// с заданной конфигурацией.
package server

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/xChygyNx/metrical/internal/server/types"
)

func middlewareLogger(h http.Handler, sugar zap.SugaredLogger) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &types.ResponseData{
			Status: 0,
			Size:   0,
		}

		lrw := types.LoggingResponseWriter{
			ResponseWriter: w,
			ResponseData:   responseData,
		}

		uri := r.RequestURI
		method := r.Method

		h.ServeHTTP(&lrw, r)

		duration := time.Since(start)
		sugar.Infoln(
			"uri", uri,
			"method", method,
			"status", responseData.Status,
			"duration", duration,
			"size", responseData.Size,
		)
	}
	return logFn
}

func getChiRouter(storage *types.MemStorage, syncInfo *types.SyncInfo,
	config *Config, sugar zap.SugaredLogger) chi.Router {
	router := chi.NewRouter()
	router.Use(GzipHandler)
	router.Post("/update",
		middlewareLogger(SaveMetricHandle(storage, syncInfo), sugar))
	router.Post("/update/",
		middlewareLogger(SaveMetricHandle(storage, syncInfo), sugar))
	router.Post("/updates",
		middlewareLogger(SaveBatchMetricHandle(storage, syncInfo), sugar))
	router.Post("/updates/",
		middlewareLogger(SaveBatchMetricHandle(storage, syncInfo), sugar))
	router.Post("/update/{mType}/{metric}/{value}",
		middlewareLogger(SaveMetricHandleOld(storage, syncInfo), sugar))
	router.Get("/value/{mType}/{metric}",
		middlewareLogger(GetMetricHandle(storage), sugar))
	router.Post("/value",
		middlewareLogger(GetJSONMetricHandle(storage), sugar))
	router.Post("/value/",
		middlewareLogger(GetJSONMetricHandle(storage), sugar))
	router.Get("/ping", middlewareLogger(pingDBHandle(config.DBAddress), sugar))
	router.Get("/", middlewareLogger(ListMetricHandle(storage), sugar))
	return router
}

func configAndSync(storage *types.MemStorage) (config *Config, syncInfo *types.SyncInfo, err error) {
	config, err = GetConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("error in GetConfig: %w", err)
	}

	if config.Restore {
		err = restoreMetricStore(config.FileStoragePath, storage)
		var storageFileNotFound *fs.PathError
		if errors.As(err, &storageFileNotFound) {

		} else if err != nil {
			return nil, nil, fmt.Errorf("error with restore MemStorage from file: %w", err)
		}
	}

	syncInfo, err = GetSyncInfo(*config)
	if err != nil {
		return nil, nil, fmt.Errorf("error in GetSyncInfo: %w", err)
	}

	return
}

// Routing запускает сервер по приему http запросов на сохранение метрик в хранилища
// указанные в конфигурации.
func Routing() (err error) {
	// Initialize logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		return errors.New("error in create zap registrator")
	}
	defer func() {
		err := logger.Sync()
		if err != nil {
			return
		}
	}()
	sugar := *logger.Sugar()
	storage := types.GetMemStorage()

	config, syncInfo, err := configAndSync(storage)
	if err != nil {
		return fmt.Errorf("error in configAndSync: %w", err)
	}

	if syncInfo.DB != nil {
		defer func() {
			err = syncInfo.DB.Close()
		}()
	} else if syncInfo.DB == nil && !syncInfo.SyncFileRecord {
		go func() {
			err = fileDump(config.FileStoragePath, time.Duration(config.StoreInterval)*time.Second, storage)
		}()
	}

	router := getChiRouter(storage, syncInfo, config, sugar)

	err = http.ListenAndServe(config.HostPort.String(), router)
	if err != nil {
		return fmt.Errorf("error with launch http server: %w", err)
	}
	return
}
