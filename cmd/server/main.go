package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sonkiee/ehhh-api/internal/handler"
	"github.com/sonkiee/ehhh-api/internal/repository"
	"github.com/sonkiee/ehhh-api/internal/routes"
	"github.com/sonkiee/ehhh-api/internal/service"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// TODO: swap these with Postgres implementations later
	dRepo := repository.NewInMemoryDilemmaRepo() // implement
	vRepo := repository.NewInMemoryVoteRepo()    // implement

	svc := service.NewDilemmaService(dRepo, vRepo)
	h := handler.NewDilemmaHandler(svc)

	r := routes.Setup(h)

	log.Printf("listening on :%s", port)
	log.Fatal(r.Run(":" + port))
}
