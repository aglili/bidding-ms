package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aglili/auction-app/internal/config"
	"github.com/aglili/auction-app/internal/provider"
	"github.com/aglili/auction-app/internal/routes"
)

func main() {
	cfg := config.LoadConfig()

	db, err := config.ConnectToDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database : %v", err)
	}
	defer db.Close()

	redis, err := config.ConnectToRedis(cfg)
	if err != nil {
		log.Fatalf("failed to connect to redis : %v", err)
	}
	defer redis.Close()

	prov := provider.NewProvider(cfg, db, redis)

	routes := routes.SetupRoutes(prov)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.AppPort),
		Handler: routes,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}