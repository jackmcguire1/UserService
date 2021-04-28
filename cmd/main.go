package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	apex "github.com/apex/log"
	"github.com/jackmcguire1/UserService/api"
	"github.com/jackmcguire1/UserService/dom/user"
)

var (
	userService user.UserService
	userHandler *api.UserHandler

	elasticSearchHost       string
	elasticSearchPort       string
	elasticSearchSecondPort string
	elasticSearchUserIndex string

	listenPort string
	listenHost string
)

func init() {
	elasticSearchHost = os.Getenv("ELASTIC_HOST")
	elasticSearchPort = os.Getenv("ELASTIC_PORT")
	elasticSearchSecondPort = os.Getenv("ELASTIC_SECOND_PORT")
	elasticSearchUserIndex = os.Getenv("ELASTIC_USER_INDEX")
	listenPort = os.Getenv("LISTEN_PORT")
	listenHost = os.Getenv("LISTEN_HOST")

	var err error
	userService, err = user.NewService(&user.Resources{
		Repo: user.NewElasticRepo(&user.ElasticSearchParams{
			Host:          elasticSearchHost,
			Port:          elasticSearchPort,
			SecondPort:    elasticSearchPort,
			UserIndexName: elasticSearchUserIndex,
		}),
	})
	if err != nil {
		log.Fatal(err)
	}

	userHandler = &api.UserHandler{UserService: userService}
}

func main() {
	s := http.NewServeMux()

	s.Handle("/user", userHandler)

	addr := fmt.Sprintf("%s:%s", listenHost, listenPort)

	apex.
		WithField("listen-address", addr).
		Info("listening")

	err := http.ListenAndServe(addr, s)
	if err != nil {
		apex.
			WithError(err).
			Fatal("failed to listen and serve")
	}
}
