package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/galrub/go/jobSearch/internal/logger"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

type logFields struct {
	ID         string
	RemoteIP   string
	Host       string
	Method     string
	Path       string
	Protocol   string
	StatusCode int
	Latency    float64
	Error      error
	Stack      []byte
}

func (lf *logFields) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("id", lf.ID).
		//		Str("remote_ip", lf.RemoteIP).
		Str("host", lf.Host).
		Str("method", lf.Method).
		Str("path", lf.Path).
		//		Str("protocol", lf.Protocol).
		Int("status_code", lf.StatusCode)
		//		Float64("latency", lf.Latency).
		//		Str("tag", "request")

	if lf.Error != nil {
		e.Err(lf.Error)
	}

	if lf.Stack != nil {
		e.Bytes("stack", lf.Stack)
	}
}

func LoggerMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rId := r.Header.Get("X-Request-ID")
		if rId == "" {
			rId = uuid.New().String()
		}

		fields := &logFields{
			ID:       rId,
			RemoteIP: r.RemoteAddr,
			Host:     r.Host,
			Method:   r.Method,
			Path:     r.URL.Path,
			Protocol: r.Proto,
		}
		lrw := NewLoggingResponseWriter(w)
		defer func() {
			rvr := recover()
			fields.StatusCode = lrw.statusCode
			if rvr != nil {
				err, ok := rvr.(error)
				if !ok {
					err = fmt.Errorf("%v", rvr)
				}
				fields.Error = err
				fields.Stack = debug.Stack()
				fields.StatusCode = http.StatusInternalServerError
			}
			fields.Latency = time.Since(start).Seconds()
			switch {
			case rvr != nil:
				logger.LOG.Error().EmbedObject(fields).Msg("Panic!")
			case fields.StatusCode > 500:
				logger.LOG.Error().EmbedObject(fields).Msg("Internal Server error")
			case fields.StatusCode > 400:
				logger.LOG.Error().EmbedObject(fields).Msg("Client error")
			case fields.StatusCode > 300:
				logger.LOG.Warn().EmbedObject(fields).Msg("Redirect")
			case fields.StatusCode > 200:
				logger.LOG.Info().EmbedObject(fields).Msg("Success")
			case fields.StatusCode > 100:
				logger.LOG.Info().EmbedObject(fields).Msg("Information")
			}
		}()
		next.ServeHTTP(lrw, r)
	})
}
