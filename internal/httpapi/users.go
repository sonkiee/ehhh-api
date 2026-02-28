package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sonkiee/ehhh-api/internal/repo"
	"github.com/sonkiee/ehhh-api/internal/util"
)

type UsersAPI struct {
	q *repo.Queries
}

func NewUsersAPI(q *repo.Queries) http.Handler {
	api := &UsersAPI{q: q}
	r := chi.NewRouter()
	r.Post("/", api.create)
	return r
}

func (a *UsersAPI) create(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Username == "" {
		util.WriteError(w, http.StatusBadRequest, "invalid body")
		// util.WriteError(w, 400, "invalid body")
		return
	}

	_, err := a.q.GetUserByUsername(r.Context(), req.Username)
	if err == nil {
		util.WriteError(w, http.StatusBadRequest, ErrUsernameTaken.Error())
		return
	}

	u, err := a.q.CreateUser(r.Context(), req.Username)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	util.WriteJSON(w, http.StatusCreated, u)
}
