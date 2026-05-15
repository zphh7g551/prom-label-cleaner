package cache

import (
	"net/http"
	"strconv"
	"time"
)

// Middleware wraps an http.Handler, serving cached responses when available
// and delegating to the upstream handler on a cache miss.
func Middleware(c *Cache, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if entry, ok := c.Get(); ok {
			age := int(time.Since(entry.FetchedAt).Seconds())
			w.Header().Set("X-Cache", "HIT")
			w.Header().Set("Age", strconv.Itoa(age))
			w.Header().Set("Content-Type", "text/plain; version=0.0.4")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(entry.Data)
			return
		}

		rec := &responseRecorder{header: make(http.Header), code: http.StatusOK}
		next.ServeHTTP(rec, r)

		if rec.code == http.StatusOK {
			c.Set(rec.body)
		}

		for k, vals := range rec.header {
			for _, v := range vals {
				w.Header().Add(k, v)
			}
		}
		w.Header().Set("X-Cache", "MISS")
		w.WriteHeader(rec.code)
		_, _ = w.Write(rec.body)
	})
}

type responseRecorder struct {
	header http.Header
	code   int
	body   []byte
}

func (r *responseRecorder) Header() http.Header        { return r.header }
func (r *responseRecorder) WriteHeader(code int)       { r.code = code }
func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body = append(r.body, b...)
	return len(b), nil
}
