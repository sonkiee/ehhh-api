package httpapi

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Deps struct {
	UserHandler    http.Handler
	DilemmaHandler http.Handler
	VoteHandler    http.Handler
	CommentHandler http.Handler
	ReportHandler  http.Handler
	Timeout        time.Duration
}

func NewRouter(d Deps) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(d.Timeout))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	r.Route("/api", func(r chi.Router) {
		r.Mount("/users", d.UserHandler)
		r.Mount("/dilemmas", d.DilemmaHandler)
		r.Mount("/votes", d.VoteHandler)
		r.Mount("/comments", d.CommentHandler)
		r.Mount("/reports", d.ReportHandler)
	})
	return r
}
