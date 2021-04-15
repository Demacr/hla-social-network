package main

import (
	"log"
	"net/http"

	"github.com/pkg/errors"

	"github.com/Demacr/otus-hl-socialnetwork/internal/config"
	"github.com/Demacr/otus-hl-socialnetwork/internal/models"
	"github.com/Demacr/otus-hl-socialnetwork/internal/storages"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	config, err := config.Configure()
	if err != nil {
		log.Fatal(err)
	}

	storage := storages.NewDB(config)

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	e.POST("/api/registrate", func(c echo.Context) error {
		profile := &models.Profile{}
		if err := c.Bind(profile); err != nil {
			return err
		}

		err := storage.WriteProfile(profile)
		if err != nil {
			log.Println("WriteProfile error:", err)
			return c.String(http.StatusInternalServerError, "")
		}

		return c.String(http.StatusCreated, "")
	})

	e.POST("/api/authorize", func(c echo.Context) error {
		credentials := &models.Credentials{}
		if err := c.Bind(credentials); err != nil {
			log.Println(err)
			return c.String(http.StatusBadRequest, "Bad json.")
		}

		result, err := storage.CheckCredentials(credentials)
		if err != nil {
			err = errors.Wrap(err, "Authorization failed")
			log.Println(err)
			return c.String(http.StatusInternalServerError, "Check credentials error.")
		}
		if !result {
			return c.String(http.StatusInternalServerError, "Wrong credentials.")
		}
		return c.String(http.StatusOK, "authorized!")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
