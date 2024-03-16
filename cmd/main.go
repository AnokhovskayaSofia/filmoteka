package main

import (
	"filmoteka/api"
	"filmoteka/config"
	"filmoteka/db"
	"filmoteka/logger"
	"log/slog"
	"net/http"
)

func main() {
	cfg := config.CnfLoad()

	log := logger.SetupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	pgdb, err := db.StartDB(cfg)
	if err != nil {
		log.Error("error starting the database %v", err)
	}

	router := api.StartAPI(pgdb)

	err = http.ListenAndServe(":8085", router)
	if err != nil {
		log.Error("error from router %v\n", err)
	}
}
