// Package httpapi contains HTTP-specific behavior: routes, multipart parsing,
// request-size limits, JSON responses, and error mapping. Handlers stay THIN —
// they translate HTTP <-> the app service and never do media work themselves.
//
// STORY 11: implement postAsset and Routes.
// STORY 12: add health endpoints and server timeouts.
// STORY 13: add logging middleware.
package httpapi

import (
	"net/http"

	"github.com/PeterShin23/go-image-processor/internal/app"
	"github.com/PeterShin23/go-image-processor/internal/domain"
)

// Handlers holds the dependencies the HTTP layer needs. Injected, not global.
type Handlers struct {
	svc       *app.Service
	maxUpload int64 // bytes; enforced before/while reading the multipart body
}

// NewHandlers wires the HTTP layer to the shared application service.
func NewHandlers(svc *app.Service, maxUpload int64) *Handlers {
	return &Handlers{svc: svc, maxUpload: maxUpload}
}

// Routes returns the HTTP handler with all routes registered.
//
// STORY 11: register POST /v1/assets -> h.postAsset.
// STORY 12: register GET /health/live and GET /health/ready.
func (h *Handlers) Routes() http.Handler {
	mux := http.NewServeMux()
	// TODO (story 11): mux.HandleFunc("POST /v1/assets", h.postAsset)
	// TODO (story 12): health routes
	return mux
}

// postAsset accepts one multipart image, calls the app service, and writes a
// JSON manifest. It must map errors to status codes and clean up temp files.
//
// TODO (story 11): parse multipart (r.ParseMultipartForm / r.FormFile),
// enforce h.maxUpload (http.MaxBytesReader), build an app.ProcessInput, call
// h.svc.Process, and encode the manifest as JSON. Map validation errors to 4xx
// and unexpected errors to 5xx.
func (h *Handlers) postAsset(w http.ResponseWriter, r *http.Request) {
	writeJSONError(w, http.StatusNotImplemented, domain.ErrNotImplemented.Error())
}

// live is the liveness probe: is the process up at all?
//
// TODO (story 12): respond 200 with a small JSON body like {"status":"ok"}.
func (h *Handlers) live(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// ready is the readiness probe: can the process serve requests (deps ready)?
//
// TODO (story 12): respond 200 when ready. Keep it cheap.
func (h *Handlers) ready(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

// writeJSON encodes v as JSON with the given status. (Small shared helper.)
//
// TODO (story 11): set Content-Type, WriteHeader(status), json.NewEncoder(w).Encode(v).
func writeJSON(w http.ResponseWriter, status int, v any) {
	// TODO
}

// writeJSONError writes a JSON error body like {"error": "..."}.
//
// TODO (story 11): implement using writeJSON.
func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
}
