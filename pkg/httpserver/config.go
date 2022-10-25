package httpserver

import (
	"context"
	"crypto/tls"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nexmoinc/gosrvlib/pkg/httpserver/route"
	"github.com/nexmoinc/gosrvlib/pkg/ipify"
	"github.com/nexmoinc/gosrvlib/pkg/profiling"
	"github.com/nexmoinc/gosrvlib/pkg/redact"
	"github.com/nexmoinc/gosrvlib/pkg/traceid"
)

// RedactFn is an alias for a redact function.
type RedactFn func(s string) string

// IndexHandlerFunc is a type alias for the route index function.
type IndexHandlerFunc func([]route.Route) http.HandlerFunc

// GetPublicIPFunc is a type alias for function to get public IP of the service.
type GetPublicIPFunc func(context.Context) (string, error)

// GetPublicIPDefaultFunc returns the GetPublicIP function for a default ipify client.
func GetPublicIPDefaultFunc() GetPublicIPFunc {
	c, _ := ipify.New() // no errors are returned with default values
	return c.GetPublicIP
}

type config struct {
	router                  Router
	serverAddr              string
	serverReadHeaderTimeout time.Duration
	serverReadTimeout       time.Duration
	serverWriteTimeout      time.Duration
	shutdownTimeout         time.Duration
	tlsConfig               *tls.Config
	defaultEnabledRoutes    []defaultRoute
	indexHandlerFunc        IndexHandlerFunc
	ipHandlerFunc           http.HandlerFunc
	metricsHandlerFunc      http.HandlerFunc
	pingHandlerFunc         http.HandlerFunc
	pprofHandlerFunc        http.HandlerFunc
	statusHandlerFunc       http.HandlerFunc
	traceIDHeaderName       string
	redactFn                RedactFn
	middleware              []MiddlewareFn
}

func defaultConfig() *config {
	return &config{
		serverAddr:              ":8017",
		serverReadHeaderTimeout: 1 * time.Minute,
		serverReadTimeout:       1 * time.Minute,
		serverWriteTimeout:      1 * time.Minute,
		shutdownTimeout:         30 * time.Second,
		defaultEnabledRoutes:    nil,
		indexHandlerFunc:        defaultIndexHandler,
		ipHandlerFunc:           defaultIPHandler(GetPublicIPDefaultFunc()),
		metricsHandlerFunc:      notImplementedHandler,
		pingHandlerFunc:         defaultPingHandler,
		pprofHandlerFunc:        profiling.PProfHandler,
		statusHandlerFunc:       defaultStatusHandler,
		traceIDHeaderName:       traceid.DefaultHeader,
		redactFn:                redact.HTTPData,
		middleware:              []MiddlewareFn{},
	}
}

func (c *config) isIndexRouteEnabled() bool {
	for _, r := range c.defaultEnabledRoutes {
		if r == IndexRoute {
			return true
		}
	}

	return false
}

// validate the configuration.
// nolint: gocyclo
func (c *config) validate() error {
	if err := validateAddr(c.serverAddr); err != nil {
		return err
	}

	if c.shutdownTimeout == 0 {
		return fmt.Errorf("invalid shutdownTimeout")
	}

	if c.ipHandlerFunc == nil {
		return fmt.Errorf("ipHandlerFunc is required")
	}

	if c.metricsHandlerFunc == nil {
		return fmt.Errorf("metricsHandlerFunc is required")
	}

	if c.pingHandlerFunc == nil {
		return fmt.Errorf("pingHandlerFunc is required")
	}

	if c.pprofHandlerFunc == nil {
		return fmt.Errorf("pprofHandlerFunc is required")
	}

	if c.statusHandlerFunc == nil {
		return fmt.Errorf("statusHandlerFunc is required")
	}

	if c.traceIDHeaderName == "" {
		return fmt.Errorf("traceIDHeaderName is required")
	}

	if c.router == nil {
		return fmt.Errorf("router is required")
	}

	return nil
}

// validateAddr checks if a http server bind address is valid.
func validateAddr(addr string) error {
	if !strings.Contains(addr, ":") {
		return fmt.Errorf("invalid http server address: %s", addr)
	}

	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid http server address: %s", addr)
	}

	port := parts[1]
	if port == "" {
		return fmt.Errorf("invalid http server address: %s", addr)
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("invalid http server address: %s", addr)
	}

	if portInt < 1 || portInt > math.MaxUint16 {
		return fmt.Errorf("invalid http server address: %s", addr)
	}

	return nil
}

func (c *config) defaultRouter(ctx context.Context) {
	r := httprouter.New()
	l := logging.FromContext(ctx)

	middleware := []MiddlewareFn{
		loggerMiddleware(l, c.traceIDHeaderName, c.redactFn),
	}
	middleware = append(middleware, c.middleware...)

	r.NotFound = applyMiddleware(
		Route{
			Path: "404",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				httputil.SendStatus(r.Context(), w, http.StatusNotFound)
			},
		},
		middleware...,
	)

	r.MethodNotAllowed = applyMiddleware(
		Route{
			Path: "405",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				httputil.SendStatus(r.Context(), w, http.StatusMethodNotAllowed)
			},
		},
		middleware...,
	)

	r.PanicHandler = func(w http.ResponseWriter, r *http.Request, p interface{}) {
		applyMiddleware(
			Route{
				Path: "500",
				Handler: func(w http.ResponseWriter, r *http.Request) {
					logging.FromContext(r.Context()).Error(
						"panic",
						zap.Any("err", p),
						zap.String("stacktrace", string(debug.Stack())),
					)
					httputil.SendStatus(r.Context(), w, http.StatusInternalServerError)
				},
			},
			middleware...,
		).ServeHTTP(w, r)
	}

	c.router = r
}
