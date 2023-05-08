package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Host          string `envconfig:"HOST" default:"0.0.0.0"`
	Port          string `envconfig:"PORT" default:"8080"`
	Domain        string `envconfig:"DOMAIN" required:"true"`
	LndHost       string `envconfig:"LND_HOST" required:"true"`
	LndMacaroon   string `envconfig:"LND_MACAROON" required:"true"`
	LndCert       string `envconfig:"LND_CERT" default:""`
	NostrKey      string `envconfig:"NOSTR_KEY" required:"true"`
	CommentLength int64  `envconfig:"COMMENT_LENGTH" default:"0"`
	WebhookUrl    string `envconfig:"WEBHOOK_URL" default:""`
	Description   string `envconfig:"DESCRIPTION" default:"Send sats!"`
	ImageUrl      string `envconfig:"IMAGE_URL" default:""`
}

var config Config
var router = mux.NewRouter()

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't process envconfig.")
	}

	router.Path("/.well-known/lnurlp/{username}").Methods("GET", "OPTIONS").HandlerFunc(handleLNURL)

	srv := &http.Server{
		Handler:      router,
		Addr:         config.Host + ":" + config.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Debug().Str("addr", srv.Addr).Msg("listening")
	srv.ListenAndServe()
}
