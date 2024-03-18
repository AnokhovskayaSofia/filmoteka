package main

import (
	"filmoteka/api"
	"filmoteka/config"
	"filmoteka/db"
	"filmoteka/logger"
	"log/slog"
	"net/http"
)

//	@title			Filmoteka API
//	@version		1.0
//	@description	This is a sample Filmoteka server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

// @host		host:port
// @BasePath /

// @securityDefinitions.basic BasicAuth
// @scope.admin Grants read and write access to administrative information
// @in header
// @name Authorization
func main() {
	cfg := config.CnfLoad()

	log := logger.SetupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	pgdb, err := db.StartDB(cfg)
	if err != nil {
		log.Error("error starting the database %v", err)
	}

	router := api.StartAPI(pgdb, cfg)

	err = http.ListenAndServe(cfg.HTTPServer.Address, router)
	if err != nil {
		log.Error("error from router %v\n", err)
	}
}
