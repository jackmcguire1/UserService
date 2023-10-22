package main

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackmcguire1/UserService/api/healthcheck"
	"github.com/jackmcguire1/UserService/api/searchapi"
	"github.com/jackmcguire1/UserService/api/userapi"
	"github.com/jackmcguire1/UserService/dom/user"
	"github.com/jackmcguire1/UserService/pkg/utils"
)

var (
	log                *slog.Logger
	userService        user.UserService
	userHandler        *userapi.UserHandler
	searchHandler      *searchapi.SearchHandler
	healthCheckHandler *healthcheck.HealthCheckHandler

	mongoHost            string
	mongoDatabase        string
	mongoUsersCollection string

	listenPort string
	listenHost string

	userUpdates chan *user.UserUpdate
	eventsURL   string
)

func init() {
	jsonLogHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	log = slog.New(jsonLogHandler)

	mongoHost = os.Getenv("MONGO_HOST")
	mongoDatabase = os.Getenv("MONGO_DATABASE")
	mongoUsersCollection = os.Getenv("MONGO_USERS_COLLECTION")

	listenPort = os.Getenv("LISTEN_PORT")
	listenHost = os.Getenv("LISTEN_HOST")

	userUpdates = make(chan *user.UserUpdate, 1)
	eventsURL = os.Getenv("EVENTS_URL")

	var err error

	userMongoRepo, err := user.NewMongoRepo(context.Background(), &user.MongoRepoParams{
		Host:           mongoHost,
		Database:       mongoDatabase,
		CollectionName: mongoUsersCollection,
	})
	if err != nil {
		log.
			With("error", err).
			Error("failed to init user mongo repo")
		panic(err)
	}

	userService, err = user.NewService(&user.Resources{
		UserChannel: userUpdates,
		Repo:        userMongoRepo,
	})
	if err != nil {
		log.
			With("error", err).
			Error("failed to init user service")
		panic(err)
	}

	userHandler = &userapi.UserHandler{UserService: userService, Logger: log}
	searchHandler = &searchapi.SearchHandler{UserService: userService, Logger: log}
	healthCheckHandler = &healthcheck.HealthCheckHandler{LogVerbosity: "DEBUG", StartTime: time.Now().UTC(), Logger: log}
}

func main() {
	slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	s := http.NewServeMux()

	s.Handle("/users", userHandler)
	s.HandleFunc("/search/users/by_country", searchHandler.UsersByCountry)
	s.HandleFunc("/search/users/", searchHandler.GetAllUsers)
	s.Handle("/healthcheck", healthCheckHandler)

	addr := fmt.Sprintf("%s:%s", listenHost, listenPort)

	log.
		With("events-url", eventsURL).
		Info("starting user updates handler")

	httpServer := http.Server{
		Addr:    addr,
		Handler: s,
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGKILL)

	go func() {
		_ = <-signals
		log.Warn("SERVER CLOSING.")

		err := httpServer.Close()
		if err != nil {
			log.
				With("error", err).
				Error("failed to close http server")
		}
	}()

	// POST user updates to URL
	go func(i chan *user.UserUpdate) {
		for update := range userUpdates {

			log.
				With("update", utils.ToJSON(update)).
				Info("got user update")

			if eventsURL != "" {
				r := bytes.NewReader([]byte(utils.ToJSON(update)))

				_, err := http.Post(eventsURL, "application/json", r)
				if err != nil {
					log.
						With("error", err).
						Error("failed to publish user update")
				}
			}
		}
	}(userUpdates)

	log.
		With("addr", addr).
		Info("starting http server")

	err := httpServer.ListenAndServe()
	if err != nil {
		log.
			With("error", err).
			Error("failed to listen and serve")
		panic(err)
	}
}
