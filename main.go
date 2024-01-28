package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/fiatjaf/eventstore/slicestore"
	"github.com/fiatjaf/khatru"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
)

type Settings struct {
	Port   string `envconfig:"PORT" default:"4777"`
	Domain string `envconfig:"DOMAIN"`

	LiveKitAPIKey    string `envconfig:"LK_API_KEY" required:"true"`
	LiveKitAPISecret string `envconfig:"LK_API_SECRET" required:"true"`
}

var (
	s   Settings
	db  = slicestore.SliceStore{}
	log = zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
)

func main() {
	err := envconfig.Process("", &s)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't process envconfig")
		return
	} else {
		if canonicalHost := os.Getenv("CANONICAL_HOST"); canonicalHost != "" {
			s.Domain = canonicalHost
		}
	}

	// expose our internal cache as a relay (mostly for debugging purposes)
	db.Init()
	defer db.Close()
	relay := khatru.NewRelay()
	relay.QueryEvents = append(relay.QueryEvents, db.QueryEvents)
	relay.DeleteEvent = append(relay.DeleteEvent, db.DeleteEvent)
	relay.RejectEvent = append(relay.RejectEvent,
		rejectEvent,
	)
	relay.OnEphemeralEvent = append(relay.OnEphemeralEvent,
		handleEphemeral,
	)

	// routes
	mux := relay.Router()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homePage().Render(r.Context(), w)
		return
	})

	log.Print("listening at http://0.0.0.0:" + s.Port)
	server := &http.Server{Addr: "0.0.0.0:" + s.Port, Handler: cors.Default().Handler(relay)}
	go func() {
		server.ListenAndServe()
		if err := server.ListenAndServe(); err != nil {
			log.Error().Err(err).Msg("")
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
	server.Close()
}
