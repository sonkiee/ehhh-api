package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sonkiee/ehhh-api/internal/config"
	"github.com/sonkiee/ehhh-api/internal/db"
	"github.com/sonkiee/ehhh-api/internal/httpapi"
	"github.com/sonkiee/ehhh-api/internal/repo"
)

func main() {

	// Load environment variables from .env file if it exists
	if os.Getenv("APP_ENV") != "production" {
		_ = godotenv.Load()
	}

	cfg := config.LoadConfig()

	ctx := context.Background()

	pool, err := db.NewPool(ctx, cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	q := repo.New(pool)

	router := httpapi.NewRouter(httpapi.Deps{
		UserHandler:    httpapi.NewUsersAPI(q),
		DilemmaHandler: httpapi.NewDilemmasAPI(pool, q),
		VoteHandler:    httpapi.NewVotesAPI(q),
		CommentHandler: httpapi.NewCommentsAPI(q),
		ReportHandler:  httpapi.NewReportsAPI(q),
		Timeout:        time.Duration(cfg.Timeout) * time.Second,
	})

	addr := ":" + cfg.Port
	log.Printf("starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))

}
