package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/apex/log"
	"github.com/jackmcguire1/UserService/api/healthcheck"
	"github.com/jackmcguire1/UserService/api/searchapi"
	"github.com/jackmcguire1/UserService/api/userapi"
	"github.com/jackmcguire1/UserService/dom/user"
	"github.com/jackmcguire1/UserService/pkg/utils"
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

	userUpdates chan *user.UserUpdate
	eventsURL   string
)

func init() {
	logLevel := os.Getenv("LOG_VERBOSITY")
	switch logLevel {
	case "":
		logLevel = "info"
		fallthrough
	default:
		log.SetLevelFromString(logLevel)
	}

	elasticSearchHost = os.Getenv("ELASTIC_HOST")
	elasticSearchPort = os.Getenv("ELASTIC_PORT")
	elasticSearchSecondPort = os.Getenv("ELASTIC_SECOND_PORT")
	elasticSearchUserIndex = os.Getenv("ELASTIC_USER_INDEX")

	listenPort = os.Getenv("LISTEN_PORT")
	listenHost = os.Getenv("LISTEN_HOST")

	userUpdates = make(chan *user.UserUpdate, 1)
	eventsURL = os.Getenv("EVENTS_URL")

	var err error
	userService, err = user.NewService(&user.Resources{
		UserChannel: userUpdates,
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
		WithField("events-url", eventsURL).
		Info("starting user updates handler")

	// POST user updates to URL
	go func(i chan *user.UserUpdate) {
		for update := range userUpdates {

			log.
				WithField("update", utils.ToJSON(update)).
				Info("got user update")

			if eventsURL != "" {
				r := bytes.NewReader([]byte(utils.ToJSON(update)))

				_, err := http.Post(eventsURL, "application/json", r)
				if err != nil {
					log.
						WithError(err).
						Error("failed to publish user update")
				}
			}
		}
	}(userUpdates)

	log.
		WithField("addr", addr).
		Info("starting http server")

	err := http.ListenAndServe(addr, s)
	if err != nil {
		log.
			WithError(err).
			Fatal("failed to listen and serve")
	}
}
