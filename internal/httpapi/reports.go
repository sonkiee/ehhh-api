package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sonkiee/ehhh-api/internal/repo"
	"github.com/sonkiee/ehhh-api/internal/util"
)

type ReportsAPI struct {
	q *repo.Queries
}

func NewReportsAPI(q *repo.Queries) http.Handler {
	api := &ReportsAPI{q: q}
	r := chi.NewRouter()
	r.Post("/dilemma", api.reportDilemma)
	r.Post("/comment", api.reportComment)
	return r
}

func (a *ReportsAPI) reportDilemma(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ReporterID string `json:"reporter_id"`
		DilemmaID  string `json:"dilemma_id"`
		Reason     string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, 400, "invalid body")
		return
	}

	reporterUUID, err := util.ParseUUID(req.ReporterID)
	if err != nil {
		util.WriteError(w, 400, "invalid reporter_id")
		return
	}

	dilemmaUUID, err := util.ParseUUID(req.DilemmaID)
	if err != nil {
		util.WriteError(w, 400, "invalid dilemma_id")
		return
	}

	out, err := a.q.CreateDilemmaReport(r.Context(), repo.CreateDilemmaReportParams{
		ReporterID: reporterUUID,
		DilemmaID:  dilemmaUUID,
		Reason:     req.Reason,
	})
	if err != nil {
		util.WriteError(w, 400, err.Error())
		return
	}
	util.WriteJSON(w, 201, out)
}

func (a *ReportsAPI) reportComment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ReporterID string `json:"reporter_id"`
		CommentID  string `json:"comment_id"`
		Reason     string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, 400, "invalid body")
		return
	}

	reporterUUID, err := util.ParseUUID(req.ReporterID)
	if err != nil {
		util.WriteError(w, 400, "invalid reporter_id")
		return
	}

	commentUUID, err := util.ParseUUID(req.CommentID)
	if err != nil {
		util.WriteError(w, 400, "invalid comment_id")
		return
	}

	out, err := a.q.CreateCommentReport(r.Context(), repo.CreateCommentReportParams{
		ReporterID: reporterUUID,
		CommentID:  commentUUID,
		Reason:     req.Reason,
	})
	if err != nil {
		util.WriteError(w, 400, err.Error())
		return
	}
	util.WriteJSON(w, 201, out)
}
