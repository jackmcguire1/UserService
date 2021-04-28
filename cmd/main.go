package main

import (
	"fmt"
	"github.com/jackmcguire1/UserService/api/healthcheck"
	"net/http"
	"os"

	"github.com/apex/log"
	"github.com/jackmcguire1/UserService/api/searchapi"
	"github.com/jackmcguire1/UserService/api/userapi"
	"github.com/jackmcguire1/UserService/dom/user"
)

var (
	userService        user.UserService
	userHandler        *userapi.UserHandler
	searchHandler      *searchapi.SearchHandler
	healthCheckHandler *healthcheck.HealthCheckHandler

	elasticSearchHost       string
	elasticSearchPort       string
	elasticSearchSecondPort string
	elasticSearchUserIndex  string

	listenPort string
	listenHost string
)

func init() {
	logLevel := os.Getenv("LOG_VERBOSITY")
	switch logLevel {
	case "":
		logLevel = "info"
	default:
		log.SetLevelFromString(logLevel)
	}

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
		log.WithError(err).Fatal("failed to init user service")
	}

	userHandler = &userapi.UserHandler{UserService: userService}
	searchHandler = &searchapi.SearchHandler{UserService: userService}
	healthCheckHandler = &healthcheck.HealthCheckHandler{LogVerbosity: logLevel}
}

func main() {
	s := http.NewServeMux()

	s.Handle("/user", userHandler)
	s.HandleFunc("/search/users/by_country", searchHandler.UsersByCountry)
	s.Handle("/healthcheck", healthCheckHandler)

	addr := fmt.Sprintf("%s:%s", listenHost, listenPort)

	log.
		WithField("listen-address", addr).
		Info("listening")

	go func () {

	}
	err := http.ListenAndServe(addr, s)
	if err != nil {
		log.
			WithError(err).
			Fatal("failed to listen and serve")
	}
}
