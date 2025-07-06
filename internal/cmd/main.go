package main

import (
	"log"
	"strconv"

	"github.com/ankitsalunkhe/url-shortner/internal/api"
	"github.com/ankitsalunkhe/url-shortner/internal/config"
	"github.com/ankitsalunkhe/url-shortner/internal/db"
	"github.com/ankitsalunkhe/url-shortner/internal/retriever"
	"github.com/ankitsalunkhe/url-shortner/internal/service"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("failed to load config", err)
	}

	rt, err := retriever.New(cfg.RtConfig)
	if err != nil {
		log.Fatal("unable to start zookeeper", err)
	}

	db, err := db.New()
	if err != nil {
		log.Fatal("failed to load config", err)
	}

	urlShornterService := service.New(db, rt)

	a := api.New(cfg.Port, cfg.BasePath, &urlShornterService)

	err = a.Start(":" + strconv.Itoa(cfg.Port))
	if err != nil {
		log.Fatal("failed to start echo server")
	}
}
