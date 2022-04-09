package service

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mabaro3009/example-architecture-go/infra/memory"
	"github.com/mabaro3009/example-architecture-go/pkg/hash"
	"github.com/mabaro3009/example-architecture-go/pkg/httpx"
	"github.com/mabaro3009/example-architecture-go/user"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

type Service struct {
	srv *http.Server
}

func NewService(conf *Config) (*Service, error) {
	dbs := &memoryDBs{
		user: memory.NewUserDB(),
	}
	q := &queries{
		user: dbs.user,
	}
	cmd := &commands{
		user: dbs.user,
	}
	svc := &services{
		userCreator: user.NewCreator(user.NewSimplePasswordValidator(user.DefaultMinLen), hash.NewBCryptHasher(bcrypt.DefaultCost), q.user, cmd.user),
	}

	router := mux.NewRouter()

	router.Methods(http.MethodGet).Path("/ping").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = httpx.WriteJSONResponse(w, http.StatusOK, "pong")
	})

	addUserRoutes(router, svc.userCreator, q.user)

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

type memoryDBs struct {
	user *memory.UserDB
}

type queries struct {
	user user.Queries
}

type commands struct {
	user user.Commands
}

type services struct {
	userCreator *user.Creator
}
