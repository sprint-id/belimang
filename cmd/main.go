package main

import (
	"context"
	"log"
	"net/http"

	"github.com/sprint-id/belimang/internal/cfg"
	"github.com/sprint-id/belimang/internal/handler"
	"github.com/sprint-id/belimang/internal/repo"
	"github.com/sprint-id/belimang/internal/service"
	"github.com/sprint-id/belimang/pkg/env"
	"github.com/sprint-id/belimang/pkg/postgre"
	"github.com/sprint-id/belimang/pkg/router"
	"github.com/sprint-id/belimang/pkg/validator"
)

func main() {
	env.LoadEnv()

	ctx := context.Background()
	router := router.NewRouter()
	conn := postgre.GetConn(ctx)
	defer conn.Close()
	validator := validator.New()

	cfg := cfg.Load()
	repo := repo.NewRepo(conn)
	service := service.NewService(repo, validator, cfg)
	handler.NewHandler(router, service, cfg)

	log.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalln("fail start server:", err)
	}
}
