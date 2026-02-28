package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sonkiee/ehhh-api/internal/repo"
	"github.com/sonkiee/ehhh-api/internal/util"
)

type VotesAPI struct {
	q *repo.Queries
}

func NewVotesAPI(q *repo.Queries) http.Handler {
	api := &VotesAPI{q: q}
	r := chi.NewRouter()
	r.Post("/", api.vote)
	return r
}

func (a *VotesAPI) vote(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID    string `json:"user_id"`
		DilemmaID string `json:"dilemma_id"`
		OptionID  string `json:"option_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, 400, "invalid body")
		return
	}

	userUUID, err := util.ParseUUID(req.UserID)
	if err != nil {
		util.WriteError(w, 400, "invalid user_id")
		return
	}

	dilemmaUUID, err := util.ParseUUID(req.DilemmaID)
	if err != nil {
		util.WriteError(w, 400, "invalid dilemma_id")
		return
	}

	optionUUID, err := util.ParseUUID(req.OptionID)
	if err != nil {
		util.WriteError(w, 400, "invalid option_id")
		return
	}

	out, err := a.q.CreateVoteTx(r.Context(), repo.CreateVoteTxParams{
		UserID:    userUUID,
		DilemmaID: dilemmaUUID,
		OptionID:  optionUUID,
	})

	if err != nil {
		// duplicate vote will cause a unique constraint violation error, which we can interpret as a conflict
		util.WriteError(w, 409, err.Error())
	}
	util.WriteJSON(w, 201, out)
}
