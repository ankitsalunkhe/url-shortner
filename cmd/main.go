package main

import (
	"log"
	"strconv"

	"github.com/ankitsalunkhe/url-shortner/api"
	"github.com/ankitsalunkhe/url-shortner/config"
	"github.com/ankitsalunkhe/url-shortner/db"
	"github.com/ankitsalunkhe/url-shortner/service"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("failed to load config", err)
	}

	db, err := db.New()
	if err != nil {
		log.Fatal("failed to load config", err)
	}

	urlShornterService := service.New(db)

	a := api.New(cfg.Port, cfg.BasePath, &urlShornterService)

	err = a.Start(":" + strconv.Itoa(cfg.Port))
	if err != nil {
		log.Fatal("failed to start echo server")
	}
}
