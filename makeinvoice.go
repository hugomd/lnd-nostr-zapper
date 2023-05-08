package main

import (
	"strconv"
	"time"

	"github.com/fiatjaf/makeinvoice"
	"github.com/tidwall/sjson"
)

func makeMetadata(username string) string {
	metadata, _ := sjson.Set("[]", "0.0", "text/identifier")
	metadata, _ = sjson.Set(metadata, "0.1", username+"@"+config.Domain)

	metadata, _ = sjson.Set(metadata, "1.0", "text/plain")
	metadata, _ = sjson.Set(metadata, "1.1", config.Description)

	if config.ImageUrl != "" {
		if b64, err := base64ImageFromURL(config.ImageUrl); err == nil {
			metadata, _ = sjson.Set(metadata, "2.0", "image/jpeg;base64")
			metadata, _ = sjson.Set(metadata, "2.1", b64)
		}
	}

	return metadata
}

func makeInvoice(username string, msat int, zapReq string) (bolt11 string, err error) {
	var description string
	if zapReq != "" {
		description = zapReq
	} else {
		description = makeMetadata(username)
	}

	backend := makeinvoice.LNDParams{
		Cert:     config.LndCert,
		Host:     config.LndHost,
		Macaroon: config.LndMacaroon,
	}

	return makeinvoice.MakeInvoice(makeinvoice.Params{
		Msatoshi:           int64(msat),
		UseDescriptionHash: true,
		Description:        description,
		Backend:            backend,

		Label: "lnd-nostr-zapper/" + strconv.FormatInt(time.Now().Unix(), 16),
	})
}
