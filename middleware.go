package main

import (
	"log"
	"net/http"
	"time"
)

type responseWriterWithStatusCode struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterWithStatusCode) writeHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		writer := &responseWriterWithStatusCode{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(writer, r)
		log.Println(writer.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}
