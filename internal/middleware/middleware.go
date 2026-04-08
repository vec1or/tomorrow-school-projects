package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func Recovery(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Error("Panic!", "Error", rec)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func Logging(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		EndTime := time.Since(startTime)
		logger.Info("Request handled", "Method", r.Method,
			"Path", r.URL.Path,
			"Duration", EndTime)
	})
}

func AllowMethods(methods []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, val := range methods {
			if r.Method == val {
				next.ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "Method Not Allowed!", http.StatusMethodNotAllowed)
	})
}
