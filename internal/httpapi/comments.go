package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sonkiee/ehhh-api/internal/repo"
	"github.com/sonkiee/ehhh-api/internal/util"
)

type CommentsAPI struct {
	q *repo.Queries
}

func NewCommentsAPI(q *repo.Queries) http.Handler {
	api := &CommentsAPI{q: q}
	r := chi.NewRouter()
	r.Post("/", api.create)
	return r
}

func (a *CommentsAPI) create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DilemmaID string  `json:"dilemma_id"`
		UserID    string  `json:"user_id"`
		Content   string  `json:"content"`
		ParentID  *string `json:"parent_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, 400, "invalid body")
		return
	}

	if req.DilemmaID == "" || req.UserID == "" || req.Content == "" {
		util.WriteError(w, 400, "missing required fields")
		return
	}

	dilemmaUUID, err := util.ParseUUID(req.DilemmaID)
	if err != nil {
		util.WriteError(w, 400, "invalid dilemma_id")
		return
	}

	userUUID, err := util.ParseUUID(req.UserID)
	if err != nil {
		util.WriteError(w, 400, "invalid user_id")
		return
	}

	parent := pgtype.UUID{Valid: false}
	if req.ParentID != nil && *req.ParentID != "" {
		parentUUID, err := util.ParseUUID(*req.ParentID)
		if err != nil {
			util.WriteError(w, 400, "invalid parent_id")
			return
		}
		parent = parentUUID
	}

	c, err := a.q.CreateComment(r.Context(), repo.CreateCommentParams{
		DilemmaID: dilemmaUUID,
		UserID:    userUUID,
		Content:   req.Content,
		ParentID:  parent,
	})
	if err != nil {
		util.WriteError(w, 400, err.Error())
		return
	}

	util.WriteJSON(w, 201, c)
}
