package proxy

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

// loggingMiddleware logs request line and headers before proxying.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body string
		if r.Body != nil {
			buf, _ := io.ReadAll(r.Body)
			body = string(buf)
			// Reset body so the next handler can still read it
			r.Body = io.NopCloser(bytes.NewReader(buf))
		}
		slog.Info("incoming request",
			"method", r.Method,
			"uri", r.URL.RequestURI(),
			"remote", r.RemoteAddr,
			"proto", r.Proto,
			"host", r.Host,
		)
		fmt.Println(body)
		for k, v := range r.Header {
			slog.Info("header", "key", k, "values", strings.Join(v, ", "))
		}

		// ----- capture response -----
		rec := &responseRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		// log response after handler finishes
		slog.Info("upstream response",
			"status", rec.status,
			"body", rec.buf.String(),
		)

	})
}

// responseRecorder wraps http.ResponseWriter to capture status, headers and body
type responseRecorder struct {
	http.ResponseWriter
	status int
	buf    bytes.Buffer
}

func (rw *responseRecorder) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseRecorder) Write(b []byte) (int, error) {
	// Copy into buffer for logging
	rw.buf.Write(b)
	return rw.ResponseWriter.Write(b)
}
