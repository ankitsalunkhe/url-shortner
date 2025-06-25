package main

import (
	"log"
	"strconv"

	"github.com/ankitsalunkhe/url-shortner/api"
	"github.com/ankitsalunkhe/url-shortner/config"
	"github.com/ankitsalunkhe/url-shortner/db"
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

	a := api.New(cfg.Port, cfg.BasePath, db)

	err = a.Start(":" + strconv.Itoa(cfg.Port))
	if err != nil {
		log.Fatal("failed to start echo server")
	}
}
