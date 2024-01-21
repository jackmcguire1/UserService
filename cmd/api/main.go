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

	"github.com/gorilla/mux"
	"github.com/jackmcguire1/UserService/api/auth"
	"github.com/jackmcguire1/UserService/api/healthcheck"
	"github.com/jackmcguire1/UserService/api/searchapi"
	"github.com/jackmcguire1/UserService/api/userapi"
	"github.com/jackmcguire1/UserService/dom/user"
	"github.com/jackmcguire1/UserService/pkg/utils"
)

type capturingResponseWriter struct {
	http.ResponseWriter
	body []byte
}

func (w *capturingResponseWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return w.ResponseWriter.Write(b)
}

var (
	log                *slog.Logger
	authHandler        *auth.Handler
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

	JWTSecret         []byte
	JWTExpiryDuration time.Duration
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

	JWTSecret = []byte(os.Getenv("JWT_SECRET"))
	JWTExpiryDuration = time.Hour

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

	authHandler = &auth.Handler{JWTSecret: JWTSecret, Expiry: JWTExpiryDuration}
	userHandler = &userapi.UserHandler{UserService: userService, Logger: log, AuthHandler: authHandler}
	searchHandler = &searchapi.SearchHandler{UserService: userService, Logger: log, AuthHandler: authHandler}
	healthCheckHandler = &healthcheck.HealthCheckHandler{LogVerbosity: "DEBUG", StartTime: time.Now().UTC(), Logger: log}
}

func main() {
	s := mux.NewRouter()

	s.HandleFunc("/sign_in", userHandler.SignIn)
	s.Handle("/users", userHandler)
	s.HandleFunc("/search/users/by_country", searchHandler.UsersByCountry)
	s.HandleFunc("/search/users/", searchHandler.GetAllUsers)
	s.Handle("/healthcheck", healthCheckHandler)

	addr := fmt.Sprintf("%s:%s", listenHost, listenPort)

	log.
		With("events-url", eventsURL).
		Info("starting user updates handler")

	// Use the headersMiddleware to set headers for all routes
	headersMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Content-Type", "application/json")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			if os.Getenv("DEBUG") == "true" {
				// Create a capturingResponseWriter based on the original ResponseWriter
				capturingWriter := &capturingResponseWriter{ResponseWriter: w}
				next.ServeHTTP(capturingWriter, r)
				log.
					With("request", utils.ToJSON(r)).
					With("raw-resp", utils.ToJSON(r)).
					With("raw-body", string(capturingWriter.body)).
					Debug("HTTP RESPONSE")
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
	s.Use(headersMiddleware)

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

	// Create a channel to receive signals
	signalChan := make(chan os.Signal, 1)

	// Notify the channel for interrupt, terminate, quit, and kill signals
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	// Start the HTTP server in a goroutine
	go func() {
		err := http.ListenAndServe(addr, s)
		if err != nil {
			// Handle the error if needed
			// For example, log the error or gracefully shut down the server
			panic(err)
		}
	}()

	receivedSig := <-signalChan
	log.With("signal", receivedSig).Warn("recieved OS SIGNAL")
}
