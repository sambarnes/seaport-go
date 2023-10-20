package logging

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// Middleware returns a request logging middleware
func Middleware(logger *zap.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			reqID := middleware.GetReqID(r.Context())
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			defer func() {
				remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
				if err != nil {
					remoteIP = r.RemoteAddr
				}
				scheme := "http"
				if r.TLS != nil {
					scheme = "https"
				}
				fields := []zap.Field{
					zap.Int("status_code", ww.Status()),
					zap.Int("bytes", ww.BytesWritten()),
					zap.Int64("duration", int64(time.Since(t1))),
					zap.String("duration_display", time.Since(t1).String()),
					zap.String("remote_ip", remoteIP),
					zap.String("proto", r.Proto),
					zap.String("method", r.Method),
				}
				if len(reqID) > 0 {
					fields = append(fields, zap.String("request_id", reqID))
				}
				logger.Info(fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI), fields...)
			}()

			h.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
