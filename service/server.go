package service

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mabaro3009/example-architecture-go/pkg/httpx"
	"net/http"
	"os"
	"time"
)

type Service struct {
	srv *http.Server
}

func NewService(conf Config) (*Service, error) {
	router := mux.NewRouter()

	router.Methods(http.MethodGet).Path("/ping").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = httpx.WriteJSONResponse(w, http.StatusOK, "pong")
	})

	srv := &http.Server{
		Handler: router,
		Addr:    conf.ListenAddress,
	}

	return &Service{srv: srv}, nil
}

func (s *Service) ListenAndServe() {
	if err := s.srv.ListenAndServe(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}

func (s *Service) Shutdown() {
	canCtx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()
	if err := s.srv.Shutdown(canCtx); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
