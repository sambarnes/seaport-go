package main

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/sambarnes/seaport-go/internal/logging"
	"github.com/sambarnes/seaport-go/internal/services/auth"
	"github.com/sambarnes/seaport-go/internal/services/fee"
	"github.com/sambarnes/seaport-go/internal/services/rfq"
	"github.com/sambarnes/seaport-go/pkg/proto/seaport/v1/seaportv1connect"
)

func main() {
	r := chi.NewRouter()
	r.Use(logging.Middleware(logging.Log))

	port := "8090"
	r.Mount(seaportv1connect.NewAuthServiceHandler(&auth.Service{}))
	r.Mount(seaportv1connect.NewFeesServiceHandler(&fee.Service{}))
	r.Mount(seaportv1connect.NewRFQServiceHandler(&rfq.Service{}))
	r.HandleFunc("/health", HealthCheckHandler)

	logStartupAsci("seaport", port)
	_ = http.ListenAndServe(":"+port, h2c.NewHandler(r, &http2.Server{}))
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "OK"}`)
}

func logStartupAsci(serviceName, port string) {
	startupLog, _ := zap.NewDevelopment()

	defer startupLog.Sync()
	sugar := startupLog.Sugar()
	begin := fmt.Sprintf("############# %s Startup Info ################", serviceName)
	sugar.Info(begin)
	sugar.Infof("Go version: %s", runtime.Version())
	sugar.Infof("Listening on port %s ...", port)
	sugar.Info(strings.Repeat("#", len(begin)))
}
