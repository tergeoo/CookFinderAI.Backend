package mdw

import (
	"CookFinder.Backend/pkg/puberr"
	"CookFinder.Backend/pkg/rest"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

var (
	DefaultMaxBodySize int64 = (1 << 20) * 10 // 10 MB
)

func JSON[T any](next func(w http.ResponseWriter, r *http.Request) (T, error)) http.HandlerFunc {
	return MetricsMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			r.Body = http.MaxBytesReader(w, r.Body, DefaultMaxBodySize)

			var body any
			body, err := next(w, r)
			if err != nil {
				handleError(w, r, err)
				return
			}

			switch o := body.(type) {
			case rest.JSONMsg[T]:
				w.WriteHeader(o.HTTPCode())
			default:
				w.WriteHeader(http.StatusOK)
			}

			err = json.NewEncoder(w).Encode(body)
			if err != nil {
				slog.Error(fmt.Sprintf("%+v", err))
			}
		},
	)
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	pubErr, err := puberr.ErrToPubErr(err)
	if err != nil {
		traceID := r.Context().Value("trace_id")
		requestID := r.Context().Value("request_id")

		slog.Error(
			"Private error",
			"err", err,
			"path", r.URL.Path,
			"method", r.Method,
			"query", r.URL.Query(),
			"trace_id", traceID,
			"request_id", requestID,
		)
	}

	w.WriteHeader(pubErr.HTTPCode)
	err = json.NewEncoder(w).Encode(pubErr)
	if err != nil {
		slog.Error(fmt.Sprintf("%+v", err))
	}
}
