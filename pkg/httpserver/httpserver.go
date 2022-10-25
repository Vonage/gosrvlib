// Package httpserver defines a default HTTP server with common routes.
package httpserver

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"

	"github.com/julienschmidt/httprouter"
	"github.com/nexmoinc/gosrvlib/pkg/httpserver/route"
	"github.com/nexmoinc/gosrvlib/pkg/httputil"
	"github.com/nexmoinc/gosrvlib/pkg/logging"
	"go.uber.org/zap"
)

// Router is the interface representing the router used by the HTTP http server.
type Router interface {
	http.Handler

	// Handler is an http.Handler wrapper.
	Handler(method, path string, handler http.Handler)
}

// Binder is an interface to allow configuring the HTTP router.
type Binder interface {
	// BindHTTP returns the routes.
	BindHTTP(ctx context.Context) []route.Route
}

// NopBinder returns a simple no-operation binder.
func NopBinder() Binder {
	return &nopBinder{}
}

type nopBinder struct{}

func (b *nopBinder) BindHTTP(_ context.Context) []route.Route { return nil }

// Start configures and start a new HTTP http server.
func Start(ctx context.Context, binder Binder, opts ...Option) error {
	l := logging.WithComponent(ctx, "httpserver")

	cfg := defaultConfig()

	for _, applyOpt := range opts {
		if err := applyOpt(cfg); err != nil {
			return err
		}
	}

	if cfg.router == nil {
		cfg.router = defaultRouter(ctx, cfg.traceIDHeaderName, cfg.redactFn, cfg.instrumentHandler)
	}

	if err := cfg.validate(); err != nil {
		return err
	}

	l.Debug("adding default routes")

	routes := newDefaultRoutes(cfg)

	l.Debug("adding service routes")

	customRoutes := binder.BindHTTP(ctx)

	// merge custom service routes with the default routes
	routes = append(routes, customRoutes...)

	for _, r := range routes {
		l.Debug("binding route", zap.String("path", r.Path))
		// Define the default middlewares and attach the custom middlewares.
		middlewares := []route.Middleware{
			loggerMiddleware(l, cfg.traceIDHeaderName, cfg.redactFn),
			instrumentHandler(r.Path, cfg.instrumentHandler),
		}
		middlewares = append(middlewares, r.Middlewares...)
		handler := applyMiddlewares(r.Handler, middlewares...)
		cfg.router.Handler(r.Method, r.Path, handler)
	}

	// attach route index if enabled
	if cfg.isIndexRouteEnabled() {
		l.Debug("enabling route index handler")
		cfg.router.Handler(http.MethodGet, indexPath, cfg.instrumentHandler(indexPath, cfg.indexHandlerFunc(routes)))
	}

	// wrap router with default middlewares
	return startServer(ctx, cfg)
}

func startServer(ctx context.Context, cfg *config) error {
	l := logging.FromContext(ctx)

	// create and start the http server
	s := &http.Server{
		Addr:              cfg.serverAddr,
		Handler:           cfg.router,
		ReadHeaderTimeout: cfg.serverReadHeaderTimeout,
		ReadTimeout:       cfg.serverReadTimeout,
		TLSConfig:         cfg.tlsConfig,
		WriteTimeout:      cfg.serverWriteTimeout,
	}

	// start HTTP listener
	var (
		ls  net.Listener
		err error
	)

	if cfg.tlsConfig == nil {
		ls, err = net.Listen("tcp", cfg.serverAddr)
	} else {
		ls, err = tls.Listen("tcp", cfg.serverAddr, cfg.tlsConfig)
	}

	if err != nil {
		return fmt.Errorf("failed creting the address listener: %w", err)
	}

	l.Info("listening for HTTP requests", zap.String("addr", cfg.serverAddr))

	go func() {
		if err := s.Serve(ls); err != nil {
			l.Error("unexpected HTTP server failure", zap.Error(err))
		}
	}()

	go func() {
		<-ctx.Done()

		l.Debug("shutting down HTTP http server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.shutdownTimeout)

		defer cancel()

		_ = s.Shutdown(shutdownCtx)

		l.Debug("HTTP server shutdown")
	}()

	return nil
}

func defaultRouter(ctx context.Context, traceIDHeaderName string, redactFn RedactFn, instrumentHandler InstrumentHandler) *httprouter.Router {
	r := httprouter.New()
	l := logging.FromContext(ctx)

	r.NotFound = RequestInjectHandler(l, traceIDHeaderName, redactFn, instrumentHandler("404", func(w http.ResponseWriter, r *http.Request) {
		httputil.SendStatus(r.Context(), w, http.StatusNotFound)
	}))

	r.MethodNotAllowed = RequestInjectHandler(l, traceIDHeaderName, redactFn, instrumentHandler("405", func(w http.ResponseWriter, r *http.Request) {
		httputil.SendStatus(r.Context(), w, http.StatusMethodNotAllowed)
	}))

	r.PanicHandler = func(w http.ResponseWriter, r *http.Request, p interface{}) {
		RequestInjectHandler(l, traceIDHeaderName, redactFn, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logging.FromContext(r.Context()).Error("panic",
				zap.Any("err", p),
				zap.String("stacktrace", string(debug.Stack())),
			)
			httputil.SendStatus(r.Context(), w, http.StatusInternalServerError)
		})).ServeHTTP(w, r)
	}

	return r
}

func defaultIndexHandler(routes []route.Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &route.Index{Routes: routes}
		httputil.SendJSON(r.Context(), w, http.StatusOK, data)
	}
}

func defaultIPHandler(fn GetPublicIPFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK

		ip, err := fn(r.Context())
		if err != nil {
			status = http.StatusFailedDependency
		}

		httputil.SendText(r.Context(), w, status, ip)
	}
}

func defaultPingHandler(w http.ResponseWriter, r *http.Request) {
	httputil.SendStatus(r.Context(), w, http.StatusOK)
}

func defaultStatusHandler(w http.ResponseWriter, r *http.Request) {
	httputil.SendStatus(r.Context(), w, http.StatusOK)
}

func notImplementedHandler(w http.ResponseWriter, r *http.Request) {
	httputil.SendStatus(r.Context(), w, http.StatusNotImplemented)
}
