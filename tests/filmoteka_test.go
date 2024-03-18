package tests

import (
	"filmoteka/api"
	"filmoteka/config"
	"filmoteka/db"
	"filmoteka/logger"
	"log/slog"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
)

var router *chi.Mux

func TestMain(m *testing.M) {
	cnf_var := os.Getenv("CONFIG_PATH")
	os.Setenv("CONFIG_PATH", "./config/test.yaml")
	defer os.Setenv("CONFIG_PATH", cnf_var)
	cfg := config.CnfLoad()

	log := logger.SetupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	pgdb, err := db.StartDB(cfg)
	if err != nil {
		log.Error("error starting the database %v", err)
	}

	router = api.StartAPI(pgdb, cfg)

	// err = http.ListenAndServe(cfg.HTTPServer.Address, router)
	// if err != nil {
	// 	log.Error("error from router %v\n", err)
	// }
	exitCode := m.Run()
	os.Exit(exitCode)
}

// func TestGetActors(t *testing.T) {

// 	method := "GET"
// 	url := "/actors"
// 	request, _ := http.NewRequest(method, url, bytes.NewBufferString(""))
// 	writer := httptest.NewRecorder()
// 	router.ServeHTTP(writer, request)
// 	assert.Equal(t, http.StatusOK, writer.Code)
// }

// func router() *gin.Engine {
// 	router := gin.Default()

// 	publicRoutes := router.Group("/auth")
// 	publicRoutes.POST("/register", Register)
// 	publicRoutes.POST("/login", Login)

// 	protectedRoutes := router.Group("/api")
// 	protectedRoutes.Use(middleware.JWTAuthMiddleware())
// 	protectedRoutes.POST("/entry", AddEntry)
// 	protectedRoutes.GET("/entry", GetAllEntries)

// 	return router
// }

// func setup() {
// 	err := godotenv.Load("../.env.test.local")
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	database.Connect()
// 	database.Database.AutoMigrate(&model.User{})
// 	database.Database.AutoMigrate(&model.Entry{})
// }

// func teardown() {
// 	migrator := database.Database.Migrator()
// 	migrator.DropTable(&model.User{})
// 	migrator.DropTable(&model.Entry{})
// }

// func makeRequest(method, url string, body interface{}, isAuthenticatedRequest bool) *httptest.ResponseRecorder {
// 	requestBody, _ := json.Marshal(body)
// 	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
// 	if isAuthenticatedRequest {
// 		request.Header.Add("Authorization", "Bearer "+bearerToken())
// 	}
// 	writer := httptest.NewRecorder()
// 	router().ServeHTTP(writer, request)
// 	return writer
// }

// func bearerToken() string {
// 	user := model.AuthenticationInput{
// 		Username: "yemiwebby",
// 		Password: "test",
// 	}

// 	writer := makeRequest("POST", "/auth/login", user, false)
// 	var response map[string]string
// 	json.Unmarshal(writer.Body.Bytes(), &response)
// 	return response["jwt"]
// }
