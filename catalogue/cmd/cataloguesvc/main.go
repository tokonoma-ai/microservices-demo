package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/microservices-demo/catalogue"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/weaveworks/common/middleware"
)

const (
	ServiceName = "catalogue"
)

var (
	HTTPLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Time (in seconds) spent serving HTTP requests.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path", "status_code", "isWS"})
	HTTPInflightRequests = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "http_requests_inflight",
		Help: "Current number of inflight HTTP requests.",
	}, []string{"method", "path"})
	HTTPRequestBodySize = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_body_size_bytes",
		Help:    "Size of HTTP request bodies.",
		Buckets: prometheus.ExponentialBuckets(128, 2, 10),
	}, []string{"method", "path"})
	HTTPResponseBodySize = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_response_body_size_bytes",
		Help:    "Size of HTTP response bodies.",
		Buckets: prometheus.ExponentialBuckets(128, 2, 10),
	}, []string{"method", "path"})
)

func init() {
	prometheus.MustRegister(HTTPLatency)
	prometheus.MustRegister(HTTPInflightRequests)
	prometheus.MustRegister(HTTPRequestBodySize)
	prometheus.MustRegister(HTTPResponseBodySize)
}

func main() {
	var (
		port   = flag.String("port", "80", "Port to bind HTTP listener") // TODO(pb): should be -addr, default ":80"
		images = flag.String("images", "./images/", "Image path")
		dsn    = flag.String("DSN", "catalogue_user:default_password@tcp(catalogue-db:3306)/socksdb", "Data Source Name: [username[:password]@][protocol[(address)]]/dbname")
		zip    = flag.String("zipkin", os.Getenv("ZIPKIN"), "Zipkin address")
	)
	flag.Parse()

	fmt.Fprintf(os.Stderr, "images: %q\n", *images)
	abs, err := filepath.Abs(*images)
	fmt.Fprintf(os.Stderr, "Abs(images): %q (%v)\n", abs, err)
	pwd, err := os.Getwd()
	fmt.Fprintf(os.Stderr, "Getwd: %q (%v)\n", pwd, err)
	files, _ := filepath.Glob(*images + "/*")
	fmt.Fprintf(os.Stderr, "ls: %q\n", files) // contains a list of all files in the current directory

	// Mechanical stuff.
	errc := make(chan error)

	// Log domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var tracer stdopentracing.Tracer
	{
		if *zip == "" {
			tracer = stdopentracing.NoopTracer{}
		} else {
			// Find service local IP.
			conn, err := net.Dial("udp", "8.8.8.8:80")
			if err != nil {
				logger.Log("err", err)
				os.Exit(1)
			}
			localAddr := conn.LocalAddr().(*net.UDPAddr)
			host := strings.Split(localAddr.String(), ":")[0]
			defer conn.Close()
			logger := log.With(logger, "tracer", "Zipkin")
			logger.Log("addr", zip)

			// Create Zipkin HTTP reporter
			reporter := zipkinhttp.NewReporter(*zip)

			// Create local endpoint
			endpoint, err := zipkin.NewEndpoint(ServiceName, fmt.Sprintf("%s:%s", host, *port))
			if err != nil {
				logger.Log("err", err)
				os.Exit(1)
			}

			// Create native Zipkin tracer
			nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
			if err != nil {
				logger.Log("err", err)
				os.Exit(1)
			}

			// Wrap with OpenTracing
			tracer = zipkinot.Wrap(nativeTracer)
		}
		stdopentracing.InitGlobalTracer(tracer)
	}

	// Data domain.
	db, err := sqlx.Open("mysql", *dsn)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}
	defer db.Close()

	// Check if DB connection can be made, only for logging purposes, should not fail/exit
	err = db.Ping()
	if err != nil {
		logger.Log("Error", "Unable to connect to Database", "DSN", dsn)
	}

	// Service domain.
	var service catalogue.Service
	{
		service = catalogue.NewCatalogueService(db, logger)
		service = catalogue.LoggingMiddleware(logger)(service)
	}

	// Endpoint domain.
	endpoints := catalogue.MakeEndpoints(service, tracer)

	// HTTP router
	router := catalogue.MakeHTTPHandler(endpoints, *images, logger, tracer)

	httpMiddleware := []middleware.Interface{
		middleware.Instrument{
			Duration:         HTTPLatency,
			InflightRequests: HTTPInflightRequests,
			RequestBodySize:  HTTPRequestBodySize,
			ResponseBodySize: HTTPResponseBodySize,
			RouteMatcher:     router,
		},
	}

	// Handler
	handler := middleware.Merge(httpMiddleware...).Wrap(router)

	// Create and launch the HTTP server.
	go func() {
		logger.Log("transport", "HTTP", "port", *port)
		errc <- http.ListenAndServe(":"+*port, handler)
	}()

	// Capture interrupts.
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("exit", <-errc)
}
