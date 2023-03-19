package main

import (
	"fmt"
	"log"

	"github.com/Demacr/otus-hl-socialnetwork/internal/config"
	"github.com/Demacr/otus-hl-socialnetwork/internal/controller"
	"github.com/Demacr/otus-hl-socialnetwork/internal/storages"
	"github.com/Demacr/otus-hl-socialnetwork/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config, err := config.Configure()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "./frontend/dist",
		HTML5:  true,
		Browse: false,
	}))

	snRepo := storages.NewMysqlSocialNetworkRepository(&config.MySQL)
	cacheRepo := storages.NewRedisCache(&config.Redis)
	feeduc := usecase.NewFeederUsecase(snRepo, cacheRepo)
	snuc := usecase.NewSocialNetworkUsecase(snRepo, cacheRepo, feeduc)

	controller.NewSocialNetworkHandler(e, snuc, feeduc, config.JWTSecret)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.Port)))
}
