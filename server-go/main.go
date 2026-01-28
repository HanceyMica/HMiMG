package main

import (
	"fmt"
	"log"

	"hmimg-server-go/internal/bootstrap"
	"hmimg-server-go/internal/config"
	"hmimg-server-go/internal/db"
	"hmimg-server-go/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if err := bootstrap.EnsureUploadDir(cfg); err != nil {
		log.Fatal(err)
	}

	database, err := db.Open(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := bootstrap.AutoMigrate(database); err != nil {
		log.Fatal(err)
	}

	if err := bootstrap.SeedDefaults(database); err != nil {
		log.Fatal(err)
	}

	r := server.NewRouter(database, cfg)
	addr := fmt.Sprintf(":%d", cfg.Port)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
