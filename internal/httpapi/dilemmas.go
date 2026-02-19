package httpapi

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
		UserID      string `json:"user_id"`
		Title       string `json:"title"`
		IsAnonymous bool   `json:"is_anonymous"`
		OptionA     string `json:"option_a"`
		OptionB     string `json:"option_b"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}

	if req.UserID == "" || req.Title == "" || req.OptionA == "" || req.OptionB == "" {
		util.WriteError(w, http.StatusBadRequest, "missing required fields")
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

	d, err := qtx.CreateDilemma(ctx, repo.CreateDilemmaParams{
		UserID:      mustUUID(req.UserID),
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
	d, err := a.q.GetDilemma(r.Context(), mustUUID(id))
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

func mustUUID(s string) [16]byte {
	u := parseUUIDOrPanic(s)
	return u
}

func parseUUIDOrPanic(s string) [16]byte {
	// keep this file small: use google/uuid
	// implemented below using a tiny wrapper
	return uuidFromString(s)
}

// separated to avoid clutter
func uuidFromString(s string) [16]byte {
	u := mustParseUUID(s)
	return u
}

func mustParseUUID(s string) [16]byte {
	u, err := parseUUID(s)
	if err != nil {
		panic(err)
	}
	return u
}

func parseUUID(s string) ([16]byte, error) {
	// local inline adapter to repo UUID type (sqlc uses [16]byte for UUID by default with pgx)
	// use google/uuid
	id, err := uuidParse(s)
	if err != nil {
		return [16]byte{}, err
	}
	return id, nil
}

func uuidParse(s string) ([16]byte, error) {
	// actual google/uuid call
	// (kept separate so you can swap UUID lib later)
	u, err := parseGoogleUUID(s)
	return u, err
}

func parseGoogleUUID(s string) ([16]byte, error) {
	// import is required at top: "github.com/google/uuid"
	u, err := uuid.Parse(s)
	if err != nil {
		return [16]byte{}, err
	}
	var out [16]byte
	copy(out[:], u[:])
	return out, nil
}

var _ = context.Background()
