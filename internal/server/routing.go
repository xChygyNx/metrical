package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/xChygyNx/metrical/internal/server/types"
)

var sugar zap.SugaredLogger

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

func Routing() error {
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
	sugar = *logger.Sugar()

	storage := types.GetMemStorage()
	router := chi.NewRouter()
	router.Post("/update",
		middlewareLogger(GzipHandler(SaveMetricHandle(storage)), sugar))
	router.Post("/update/",
		middlewareLogger(GzipHandler(SaveMetricHandle(storage)), sugar))
	router.Post("/update/{mType}/{metric}/{value}",
		middlewareLogger(GzipHandler(SaveMetricHandleOld(storage)), sugar))
	router.Get("/value/{mType}/{metric}",
		middlewareLogger(GzipHandler(GetMetricHandle(storage)), sugar))
	router.Post("/value",
		middlewareLogger(GzipHandler(GetJSONMetricHandle(storage)), sugar))
	router.Post("/value/",
		middlewareLogger(GzipHandler(GetJSONMetricHandle(storage)), sugar))
	router.Get("/", middlewareLogger(ListMetricHandle(storage), sugar))

	config, err := GetConfig()
	if err != nil {
		return fmt.Errorf("error in GetConfig: %w", err)
	}

	err = http.ListenAndServe(config.HostPort.String(), router)
	if err != nil {
		return fmt.Errorf("error with launch http server: %w", err)
	}
	return nil
}
