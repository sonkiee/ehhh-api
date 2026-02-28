package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sonkiee/ehhh-api/internal/repo"
	"github.com/sonkiee/ehhh-api/internal/util"
)

type DilemmasAPI struct {
	pool *pgxpool.Pool
	q    *repo.Queries
}

func NewDilemmasAPI(pool *pgxpool.Pool, q *repo.Queries) http.Handler {
	api := &DilemmasAPI{pool: pool, q: q}
	r := chi.NewRouter()
	r.Post("/", api.create)
	r.Get("/feed", api.feed)
	r.Get("/{id}", api.get)
	return r
}

func (a *DilemmasAPI) create(w http.ResponseWriter, r *http.Request) {

	var req struct {
		UserID      string `json:"userId"`
		Title       string `json:"title"`
		IsAnonymous bool   `json:"isAnonymous"`
		OptionA     string `json:"optionA"`
		OptionB     string `json:"optionB"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}

	fmt.Printf("Received create dilemma request: %+v\n", req)

	errors := make(map[string]string)

	if req.UserID == "" {
		errors["userId"] = "userId is required"
	}
	if req.Title == "" {
		errors["title"] = "title is required"
	}
	if req.OptionA == "" {
		errors["optionA"] = "optionA is required"
	}
	if req.OptionB == "" {
		errors["optionB"] = "optionB is required"
	}

	if len(errors) > 0 {
		util.WriteJSON(w, http.StatusBadRequest, errors)
		return
	}

	ctx := r.Context()
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback(ctx)

	qtx := a.q.WithTx(tx)

	userID, err := util.ParseUUID(req.UserID)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "userId must be a valid UUID")
		return
	}
	d, err := qtx.CreateDilemma(ctx, repo.CreateDilemmaParams{
		UserID:      userID,
		Title:       req.Title,
		IsAnonymous: req.IsAnonymous,
	})
	if err != nil {
		util.WriteError(w, 400, err.Error())
		return
	}

	opt1, err := qtx.CreateDilemmaOption(ctx, repo.CreateDilemmaOptionParams{
		DilemmaID: d.ID,
		Label:     req.OptionA,
	})

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	opt2, err := qtx.CreateDilemmaOption(ctx, repo.CreateDilemmaOptionParams{
		DilemmaID: d.ID,
		Label:     req.OptionB,
	})

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := tx.Commit(ctx); err != nil {
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusCreated, map[string]any{
		"dilemma": d,
		"options": []any{opt1, opt2},
	})
}

func (a *DilemmasAPI) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	parsedID, err := util.ParseUUID(id)
	if err != nil {
		util.WriteError(w, 400, "invalid id format")
		return
	}
	d, err := a.q.GetDilemma(r.Context(), parsedID)
	if err != nil {
		util.WriteError(w, 404, "not found")
		return
	}
	opts, err := a.q.GetDilemmaOptions(r.Context(), d.ID)
	if err != nil {
		util.WriteError(w, 500, err.Error())
		return
	}
	util.WriteJSON(w, 200, map[string]any{"dilemma": d, "options": opts})
}

func (a *DilemmasAPI) feed(w http.ResponseWriter, r *http.Request) {
	limit := int32(20)
	offset := int32(0)

	rows, err := a.q.ListFeed(r.Context(), repo.ListFeedParams{Limit: limit, Offset: offset})
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	util.WriteJSON(w, http.StatusOK, rows)
}

var _ = context.Background()
