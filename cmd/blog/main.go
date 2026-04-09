package main

import (
	"blog/internal/handlers"
	"blog/internal/storage"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	godotenv.Load()

	router := handlers.NewRouter()
	router.Logger = logrus.New()

	file, err := os.OpenFile("data.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		router.Logger.Fatal(err)
	}
	defer file.Close()

	router.Logger.SetOutput(file)

	DBURL := os.Getenv("DBURL")
	if err = storage.RunMigrations(DBURL, ""); err != nil {
		router.Logger.Fatal(err)
	}

	credentials := os.Getenv("credentials")
	if credentials == "" {
		router.Logger.Fatal(err)
	}

	if err = router.DBInit(credentials); err != nil {
		router.Logger.Fatal(err)
	}

	if router.PagesPath = os.Getenv("pagesPath"); router.PagesPath == "" {
		router.Logger.Fatal(err)
	}

	router.Logger.Info("Run server ...")
	router.Logger.Fatal(http.ListenAndServe("localhost:8080", router.SetRouter()))
}
