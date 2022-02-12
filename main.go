package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/hrabalvojta/micro-dvdrental/config"
	"github.com/hrabalvojta/micro-dvdrental/films"
	"github.com/hrabalvojta/micro-dvdrental/health"
	hb "github.com/hrabalvojta/micro-dvdrental/logger"
	"github.com/hrabalvojta/micro-dvdrental/psql"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// serializedLogger is our "global" application logger
type serializedLogger struct {
	log.Logger
}

// aConfig is the application configuration object
var aConfig *config.Config

func main() {
	c := config.Init()
	setConfig(c)

	// HTTP listener configuration
	var (
		port     = c.Env.HTTPPort
		httpAddr = flag.String("http.addr", ":"+port, "HTTP Listen Address")
	)

	flag.Parse()

	// Create and configure the logger
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = hb.NewHerbertFormatLogger(logger, c.Env.LogPath, c.LogLevel())
	logger = &serializedLogger{Logger: logger}
	logger = log.With(logger,
		"context_environment", c.Env.ApplicationEnvironment,
		"timestamp", log.DefaultTimestampUTC,
	)

	// Repository initialization
	var (
		userRepo films.Repository
		err      error
	)

	fieldKeys := []string{"method"}

	userRepo, err = psql.NewPsqlUserRepository(
		psql.DefaultHost,
		psql.DefaultPort,
		psql.DefaultDatabase,
		psql.DefaultDBUser,
		psql.DefaultPassword,
		psql.DefaultSSLMode,
		psql.DefaultTimeZone,
	)
	if err != nil {
		os.Exit(1)
	}

	// Initialize the users service and wrap it with our middlewares
	us := films.NewService(userRepo)
	us = films.NewLoggingService(log.With(logger, "context_component", "users"), us)
	us = films.NewInstrumentingService(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "user_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "user_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys),
		us,
	)

	// Build and initialize our application HTTP handlers and error channels
	httpLogger := log.With(logger, "context_component", "http")
	mux := http.NewServeMux()

	mux.Handle("/api/v1/users", films.MakeHandler(us, httpLogger))
	mux.Handle("/api/v1/users/", films.MakeHandler(us, httpLogger))
	mux.Handle("/api/v1/health", health.MakeHandler(httpLogger))

	http.Handle("/", accessControl(mux))
	http.Handle("/metrics", promhttp.Handler())

	srv := http.Server{
		WriteTimeout: 300 * time.Second,
		ReadTimeout:  300 * time.Second,
		Addr:         *httpAddr,
	}

	// Define the atreides logging channels
	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", httpAddr, "message", "listening")
		socket, _ := net.Listen("tcp4", *httpAddr)
		errs <- srv.Serve(socket)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}

func setConfig(c *config.Config) {
	aConfig = c
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
