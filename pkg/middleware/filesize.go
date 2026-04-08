package middleware

import (
	"fmt"
	"net/http"
	"strconv"
)

// FileSizeMiddleware validates request body size
func FileSizeMiddleware(maxSize int64) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// For multipart uploads, Content-Length may not be accurate
			// We'll validate after reading the form

			// For non-multipart requests, check Content-Length
			if cl := r.Header.Get("Content-Length"); cl != "" {
				size, _ := strconv.ParseInt(cl, 10, 64)
				if size > maxSize {
					http.Error(w, "File too large", http.StatusRequestEntityTooLarge)
					return
				}
			}

			// For multipart, we'll validate in the handler
			// Pass through with a custom writer that tracks size
			w = &sizeLimitWriter{ResponseWriter: w, limit: maxSize}
			handler.ServeHTTP(w, r)
		})
	}
}

// sizeLimitWriter wraps ResponseWriter to track written size
type sizeLimitWriter struct {
	http.ResponseWriter
	limit int64
	size  int64
}

func (w *sizeLimitWriter) Write(b []byte) (int, error) {
	if int64(len(b)) > w.limit-w.size {
		return 0, fmt.Errorf("file too large")
	}
	w.size += int64(len(b))
	return w.ResponseWriter.Write(b)
}
