package api

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net"
	"net/http"
	"ww/internal/conf"
)

type Server struct {
	router  *mux.Router
	logger  *zap.SugaredLogger
	server  *http.Server
	grStore GRStore
}

func New(config *conf.Config, grStore GRStore) (*Server, error) {

	r := mux.NewRouter()
	r.Use(RequestID)

	r.Use(loggerHTTPMiddlewareDefault(config.LoggerDisabledHttp))
	r.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	/*
		// CORS Config
		r.Use(cors.New(cors.Options{
			AllowedOrigins:   config.Strings("server.cors.allowed_origins"),
			AllowedMethods:   config.Strings("server.cors.allowed_methods"),
			AllowedHeaders:   config.Strings("server.cors.allowed_headers"),
			AllowCredentials: config.Bool("server.cors.allowed_credentials"),
			MaxAge:           config.Int("server.cors.max_age"),
		}).Handler) */

	s := &Server{
		logger:  zap.S().With("package", "server"),
		router:  r,
		grStore: grStore,
	}

	return s, nil

}

// ListenAndServe will listen for requests
func (s *Server) ListenAndServe(config *conf.Config) error {

	s.server = &http.Server{
		Addr:    net.JoinHostPort(config.ServerHost, config.ServerPort),
		Handler: s.router,
	}
	s.logger.Infow(s.server.Addr)

	// Listen

	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		fmt.Errorf("Could not listen on %s: %v", s.server.Addr, err)
		return fmt.Errorf("Could not listen on %s: %v", s.server.Addr, err)
	}

	go func() {
		if err = s.server.Serve(listener); err != nil {
			s.logger.Fatalw("API Listen error", "error", err, "address", s.server.Addr)
		}
	}()
	s.logger.Infow("API Listening", "address", s.server.Addr)

	// s.server.ListenAndServe()

	return nil

}

// Router returns the router
func (s *Server) Router() *mux.Router {
	return s.router
}
