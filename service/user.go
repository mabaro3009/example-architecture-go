package service

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mabaro3009/example-architecture-go/pkg/httpx"
	"github.com/mabaro3009/example-architecture-go/user"
)

func addUserRoutes(router *mux.Router, creator Creator, query user.Queries) {
	router.Methods(http.MethodPost).Path("/users").HandlerFunc(handleUserCreate(creator))
	router.Methods(http.MethodGet).Path("/users/{id}").HandlerFunc(handleUserGet(query))
}

type Creator interface {
	Create(ctx context.Context, params user.CreateParams) (*user.User, error)
}

func handleUserCreate(creator Creator) http.HandlerFunc {
	type userCreateRequest struct {
		ID       string `json:"id,required"`
		Username string `json:"username,required"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	type userCreateResponse struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Role     string `json:"role"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req userCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			body := map[string]string{"error": err.Error()}
			_ = httpx.WriteJSONResponse(w, http.StatusBadRequest, body)
			return
		}

		params := user.CreateParams{
			ID:       req.ID,
			Username: req.Username,
			Password: req.Password,
			Role:     req.Role,
		}

		u, err := creator.Create(context.Background(), params)
		if err != nil {
			body := map[string]string{"error": err.Error()}
			switch err {
			case user.ErrInvalidRole, user.ErrInvalidUsername, user.ErrPasswordTooSmall:
				_ = httpx.WriteJSONResponse(w, http.StatusBadRequest, body)
			default:
				_ = httpx.WriteJSONResponse(w, http.StatusInternalServerError, body)
			}
			return
		}

		resp := userCreateResponse{
			ID:       u.ID,
			Username: u.Username,
			Role:     u.Role.String(),
		}

		_ = httpx.WriteJSONResponse(w, http.StatusCreated, resp)
	}
}

func handleUserGet(q user.GetByID) http.HandlerFunc {
	type userGetResponse struct {
		ID        string     `json:"id"`
		Username  string     `json:"username"`
		Role      string     `json:"role"`
		CreatedAt time.Time  `json:"created_at"`
		DeletedAt *time.Time `json:"deleted_at"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, ok := params["id"]
		if !ok {
			body := map[string]string{"error": "missing id in url"}
			_ = httpx.WriteJSONResponse(w, http.StatusBadRequest, body)
			return
		}

		u, err := q.GetByID(context.Background(), id)
		if err != nil {
			body := map[string]string{"error": err.Error()}
			switch err {
			case user.ErrDoesNotExist:
				_ = httpx.WriteJSONResponse(w, http.StatusNotFound, body)
			default:
				_ = httpx.WriteJSONResponse(w, http.StatusInternalServerError, body)
			}
			return
		}

		resp := userGetResponse{
			ID:        u.ID,
			Username:  u.Username,
			Role:      u.Role.String(),
			CreatedAt: u.CreatedAt,
			DeletedAt: nil,
		}

		_ = httpx.WriteJSONResponse(w, http.StatusOK, resp)
	}
}
