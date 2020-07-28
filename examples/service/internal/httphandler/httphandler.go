//go:generate mockgen -package mocks -destination ../mocks/httphandler_mocks.go . Service

package httphandler

import (
	"context"
	"net/http"

	"github.com/nexmoinc/gosrvlib/pkg/httpserver/route"
	"github.com/nexmoinc/gosrvlib/pkg/httputil"
	"github.com/nexmoinc/gosrvlib/pkg/uid"
)

// NOTE: This is a sample Service interface. It is meant to demonstrate where the business logic of a service should
// reside. Also, it adds the capability of mocking the HTTP Handler independently from the rest of the code
type Service interface {
	// NOTE: Add service functions here
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
			Path:        "/uid",
			Handler:     h.handleGenUID,
			Description: "Generates a random UID.",
		},
	}
}

func (h *HTTPHandler) handleGenUID(w http.ResponseWriter, r *http.Request) {
	httputil.SendJSON(r.Context(), w, http.StatusOK, uid.NewID64())
}
