package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"os"

	_ "github.com/SawitProRecruitment/UserService/docs"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"
	_ "github.com/joho/godotenv/autoload"
)

// @title SawitPro Swagger Example API
// @version 1.0
// @description This is a sample server Cellar server.
// @termsOfService http://swagger.io/terms/
// @license.name Apache 2.0
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @host localhost:1323
func main() {
	e := echo.New()

	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)

	// Endpoint for serving Swagger JSON
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	port := os.Getenv("APP_PORT")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}

func newServer() *handler.Server {
	dbDsn := os.Getenv("DATABASE_URL")
	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})
	opts := handler.NewServerOptions{
		Repository: repo,
	}
	return handler.NewServer(opts)
}
