package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/fiatjaf/go-lnurl"
	"github.com/gorilla/mux"
	nostr "github.com/nbd-wtf/go-nostr"
	decodepay "github.com/nbd-wtf/ln-decodepay"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/sjson"
)

func handleLNURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	username := mux.Vars(r)["username"]
	nostrPubKey, err := nostr.GetPublicKey(config.NostrKey)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get public key")
	}

	log.Info().Str("username", username).Msg("Handling LNURL")
	if amount := r.URL.Query().Get("amount"); amount == "" {
		json.NewEncoder(w).Encode(lnurl.LNURLPayParams{
			LNURLResponse:   lnurl.LNURLResponse{Status: "OK"},
			Callback:        fmt.Sprintf("https://%s/.well-known/lnurlp/%s", r.Host, username),
			MinSendable:     1000,
			MaxSendable:     100000000,
			EncodedMetadata: makeMetadata(username),
			CommentAllowed:  config.CommentLength,
			Tag:             "payRequest",
			AllowsNostr:     true,
			NostrPubkey:     nostrPubKey,
		})
	} else {
		msat, err := strconv.Atoi(amount)
		if err != nil {
			json.NewEncoder(w).Encode(lnurl.ErrorResponse("amount is not integer"))
			return
		}

		zapReqStr, _ := url.QueryUnescape(r.URL.Query().Get("nostr"))

		var zapReq nostr.Event
		if err := json.Unmarshal([]byte(zapReqStr), &zapReq); err != nil {
			log.Warn().Err(err).Msg("Failed to unmarshal zap request")
			return
		}
		valid, err := zapReq.CheckSignature()
		if !valid {
			log.Info().Msg("Zap request signature invalid")
			return
		}

		log.Info().Interface("zap request", zapReq).Msg("Parsed zap request")

		bolt11, err := makeInvoice(username, msat, zapReq.String())
		if err != nil {
			json.NewEncoder(w).Encode(
				lnurl.ErrorResponse("failed to create invoice: " + err.Error()))
			return
		}

		json.NewEncoder(w).Encode(lnurl.LNURLPayValues{
			LNURLResponse: lnurl.LNURLResponse{Status: "OK"},
			PR:            bolt11,
			Routes:        make([][]interface{}, 0),
			Disposable:    lnurl.FALSE,
			SuccessAction: lnurl.Action("Payment received!", ""),
		})

		go func() {
			inv, err := decodepay.Decodepay(bolt11)
			if err != nil {
				return
			}
			WaitForZap(inv.PaymentHash, zapReq)
			// send webhook
			go sendWebhook(bolt11, r.URL.Query().Get("comment"), msat)
		}()
	}
}

func sendWebhook(bolt11, comment string, msat int) {
	body, _ := sjson.Set("{}", "pr", bolt11)
	body, _ = sjson.Set(body, "amount", msat)
	if comment != "" {
		body, _ = sjson.Set(body, "comment", comment)
	}

	(&http.Client{Timeout: 5 * time.Second}).
		Post(config.WebhookUrl, "application/json", bytes.NewBufferString(body))
}
