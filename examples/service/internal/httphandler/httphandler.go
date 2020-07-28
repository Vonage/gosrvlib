//go:generate mockgen -package mocks -destination ../mocks/httphandler_mocks.go . Service

package httphandler

import (
	"context"
	"net/http"

	"github.com/nexmoinc/gosrvlib/pkg/httpserver/route"
	"github.com/nexmoinc/gosrvlib/pkg/httputil"
)

// NOTE: Service is the interface defining the service functions required by the handler
type Service interface {
	SayHello(ctx context.Context, name string) (string, error)
	SayHelloAgain(ctx context.Context, name string) (string, error)
}

// New creates a new instance of the HTTP handler
func New(s Service) *HTTPHandler {
	return &HTTPHandler{
		service: s,
	}
}

// HTTPHandler is the struct containing all the http handlers
type HTTPHandler struct {
	service Service
}

func (h *HTTPHandler) BindHTTP(_ context.Context) []route.Route {
	return []route.Route{
		{
			Method:      http.MethodGet,
			Path:        "/hello",
			Handler:     h.handleSayHello,
			Description: "Say hello.",
		},
	}
}

func (h *HTTPHandler) handleSayHello(w http.ResponseWriter, r *http.Request) {
	// msg, err := h.service.SayHello(r.Context(), "client")
	// if err != nil {
	// 	// NOTE: Not the right thing
	// 	httputil.SendError(w, http.StatusInternalServerError, err)
	// }

	httputil.SendJSON(r.Context(), w, http.StatusOK, "Hello!")
}
