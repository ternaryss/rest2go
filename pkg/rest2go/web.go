package rest2go

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ternaryss/rest2go/pkg/rest2go/settings"
)

type server struct {
	conf        settings.Server
	router      *http.ServeMux
	middlewares []Middleware
}

func NewServer(conf settings.Server, router *http.ServeMux, middlewares ...Middleware) *server {
	if router == nil {
		router = http.NewServeMux()
	}

	if conf.HealthCheck {
		router.HandleFunc("GET /health", healthCheck)
	}

	if conf.NotFoundHandler {
		router.HandleFunc("/", HandleNotFoundError)
	}

	return &server{
		conf:        conf,
		router:      router,
		middlewares: middlewares,
	}
}

func (s *server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.conf.Host, s.conf.Port)
	var handler http.Handler

	if len(s.middlewares) > 0 {
		handler = Middlewares(s.middlewares...)(s.router)
	} else {
		handler = s.router
	}

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	slog.Info("HTTP server started", "addr", addr, "middlewares", len(s.middlewares))

	return server.ListenAndServe()
}
