package bulkethhandler

import (
	"net/http"
	"time"
)

type requestLogHandler struct {
	handler http.Handler
}

func (h requestLogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := requestLogger(r).With(
		"method", r.Method,
		"url", r.URL,
	)
	logger.Info("Started")

	t0 := time.Now()
	h.handler.ServeHTTP(w, r)
	runtime := time.Now().Sub(t0)

	logger.InfoContext(r.Context(), "Finished", "duration", runtime)
}
