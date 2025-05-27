// Package server отвечает за прием метрик от агента и их хранения в соответствии
// с заданной конфигурацией.
package server

import (
	"errors"
	"expvar"
	"fmt"
	"io/fs"
	"net/http"
	"net/http/pprof"
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

var epoch = time.Unix(0, 0).UTC().Format(http.TimeFormat)

var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, no-store, no-transform, must-revalidate, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

var etagHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-Match",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

// NoCache is a simple piece of middleware that sets a number of HTTP headers to prevent
// a router (or subrouter) from being cached by an upstream proxy and/or client.
//
// As per http://wiki.nginx.org/HttpProxyModule - NoCache sets:
//
//	Expires: Thu, 01 Jan 1970 00:00:00 UTC
//	Cache-Control: no-cache, private, max-age=0
//	X-Accel-Expires: 0
//	Pragma: no-cache (for HTTP/1.0 proxies/clients)
func NoCache(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		// Delete any ETag headers that may have been set
		for _, v := range etagHeaders {
			if r.Header.Get(v) != "" {
				r.Header.Del(v)
			}
		}

		// Set our NoCache headers
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func Profiler() http.Handler {
	r := chi.NewRouter()
	r.Use(NoCache)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.RequestURI+"/pprof/", http.StatusMovedPermanently)
	})
	r.HandleFunc("/pprof", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.RequestURI+"/", http.StatusMovedPermanently)
	})

	r.HandleFunc("/pprof/*", pprof.Index)
	r.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/pprof/profile", pprof.Profile)
	r.HandleFunc("/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/pprof/trace", pprof.Trace)
	r.Handle("/vars", expvar.Handler())

	r.Handle("/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/pprof/threadcreate", pprof.Handler("threadcreate"))
	r.Handle("/pprof/mutex", pprof.Handler("mutex"))
	r.Handle("/pprof/heap", pprof.Handler("heap"))
	r.Handle("/pprof/block", pprof.Handler("block"))
	r.Handle("/pprof/allocs", pprof.Handler("allocs"))

	return r
}

func getChiRouter(storage *types.MemStorage, syncInfo *types.SyncInfo,
	config *Config, sugar zap.SugaredLogger) chi.Router {
	router := chi.NewRouter()
	router.Use(GzipHandler)
	router.Mount("/debug", Profiler())

	//router.Get("/", func(w http.ResponseWriter, r *http.Request) {
	//	http.Redirect(w, r, r.RequestURI+"/pprof/", http.StatusMovedPermanently)
	//})
	//router.HandleFunc("/pprof", func(w http.ResponseWriter, r *http.Request) {
	//	http.Redirect(w, r, r.RequestURI+"/", http.StatusMovedPermanently)
	//})
	//
	//router.HandleFunc("/pprof/*", pprof.Index)
	//router.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	//router.HandleFunc("/pprof/profile", pprof.Profile)
	//router.HandleFunc("/pprof/symbol", pprof.Symbol)
	//router.HandleFunc("/pprof/trace", pprof.Trace)
	//router.Handle("/vars", expvar.Handler())
	//
	//router.Handle("/pprof/goroutine", pprof.Handler("goroutine"))
	//router.Handle("/pprof/threadcreate", pprof.Handler("threadcreate"))
	//router.Handle("/pprof/mutex", pprof.Handler("mutex"))
	//router.Handle("/pprof/heap", pprof.Handler("heap"))
	//router.Handle("/pprof/block", pprof.Handler("block"))
	//router.Handle("/pprof/allocs", pprof.Handler("allocs"))

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
