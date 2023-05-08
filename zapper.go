package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	b64 "encoding/base64"
	"encoding/hex"
	"github.com/gorilla/websocket"
	nostr "github.com/nbd-wtf/go-nostr"
	"github.com/rs/zerolog/log"
)

type LNDResponse struct {
	Result Invoice `json:"result"`
}

type Invoice struct {
	Memo           string `json:"memo"`
	State          string `json:"state"`
	SettleDate     int64  `json:"settle_date,string"`
	CreationDate   int64  `json:"creation_date,string"`
	PaymentRequest string `json:"payment_request"`
	PreImage       string `json:"r_preimage"`
}

func WaitForZap(r_hash string, zapReq nostr.Event) {
	log.Info().Str("r_hash", r_hash).Msg("Waiting for Zap!")

	publicKey, err := nostr.GetPublicKey(config.NostrKey)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get public key")
	}

	r_hash_bytes, err := hex.DecodeString(r_hash)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to decode r_hash")
	}

	uEnc := b64.URLEncoding.EncodeToString(r_hash_bytes)
	formatted := fmt.Sprintf("%s/v2/invoices/subscribe/%s?method=GET", strings.Replace(config.LndHost, "https", "wss", 1), uEnc)
	authHeader := http.Header{
		"Grpc-Metadata-Macaroon": []string{config.LndMacaroon},
	}
	conn, _, err := websocket.DefaultDialer.Dial(formatted, authHeader)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to dial")
	}

	log.Info().Msg("Waiting for responses on websocket")
	for {
		var response LNDResponse
		err := conn.ReadJSON(&response)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Info().Err(err).Msg("Unexpected websocket error")
			}
			break
		}
		log.Info().Interface("response", response).Msg("Got response from LND")

		go func() {
			if response.Result.State == "SETTLED" {
				zapReceipt := makeZapReceipt(config.NostrKey, publicKey, response.Result, zapReq)
				log.Info().Interface("zapReceipt", zapReceipt).Msg("Built zap receipt")

				relays := zapReq.Tags.GetAll([]string{"relays"})[0][1:]
				relaysWithDefaults := append(relays, config.Relays...)

				for _, url := range relaysWithDefaults {
					go publish(url, zapReceipt)
				}
			}
		}()
	}
}

func publish(url string, zapReceipt nostr.Event) {
	log.Info().Str("relay", url).Msg("Connecting to relay")
	relay, err := nostr.RelayConnect(context.Background(), url)
	if err != nil {
		log.Info().Err(err).Msg("Failed to connect to relay")
		return
	}
	if _, err := relay.Publish(context.Background(), zapReceipt); err != nil {
		log.Info().Err(err).Msg("Failed to publish event to relay")
		return
	}
	log.Info().Str("relay", url).Str("ID", zapReceipt.GetID()).Msg("Published to relay")
}

func makeZapReceipt(privateKey, publicKey string, invoice Invoice, zapReq nostr.Event) nostr.Event {
	preimageBytes, _ := b64.StdEncoding.DecodeString(invoice.PreImage)
	preimageHex := hex.EncodeToString(preimageBytes)

	event := nostr.Event{
		PubKey:    publicKey,
		CreatedAt: nostr.Timestamp(invoice.SettleDate),
		Kind:      nostr.KindZap,
		Tags: nostr.Tags{
			*zapReq.Tags.GetFirst([]string{"p"}),
			*zapReq.Tags.GetFirst([]string{"e"}),
			nostr.Tag{"bolt11", invoice.PaymentRequest},
			nostr.Tag{"description", zapReq.String()},
			nostr.Tag{"preimage", preimageHex},
		},
	}

	event.Sign(privateKey)

	return event
}
