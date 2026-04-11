package server

import(
	"net/http"
)

func (h *Handler) serverError(w http.ResponseWriter, r *http.Request, err error) {
	h.App.Logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func (h *Handler) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (h *Handler) notFound(w http.ResponseWriter) {
	h.clientError(w, http.StatusNotFound)
}