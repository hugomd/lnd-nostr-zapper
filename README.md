# lnd-nostr-zapper

A server which implements [LUD-06](https://github.com/lnurl/luds/blob/luds/06.md) (LNURLp) and
[NIP-57](https://github.com/nostr-protocol/nips/blob/master/57.md).

Heavily inspired by [fiatjaf/bridgeaddr](https://github.com/fiatjaf/bridgeaddr/tree/master) and [cln-nostr-zapper](https://github.com/jb55/cln-nostr-zapper/tree/master).

# Installation

A Docker container image is available via GitHub Container Registry:
```
docker pull ghcr.io/hugomd/lnd-nostr-zapper:c1225c57
```

# Configuration
Configuration is done via environment variables:

| Environment Variable  | Required | Default value | Description |
| --------------------- | -------- | ------------- | ----------- |
| `HOST`                | false    | `0.0.0.0`      | The host to bind the HTTP server to. | 
| `PORT`                | false    | `8080`         | The port to bind the HTTP server to. | 
| `DOMAIN`              | true     | N/A            | The domain associated with the server, used in the LNURL callback. | 
| `LND_HOST`            | true     | N/A            | URL pointing to LND.                  | 
| `LND_MACAROON`        | true     | N/A            | An invoice read/write macaroon for auth with LND. | 
| `LND_CERT`            | false    | `""`           | Optional self-signed certificate to call LND. | 
| `NOSTR_KEY`           | true     | N/A            | Nostr private key, used to publish zap receipts. | 
| `COMMENT_LENGTH`      | false    | `0`            | Maximum length of associated comments, sent via webhook. | 
| `WEBHOOK_URL`         | false    | `""`           | URL to call after successful payment. | 
| `DESCRIPTION`         | false    | `"Send sats!"` | Description shown in LNURL metadata. | 
| `IMAGE_URL`           | false    | `""`           | Optional avatar URL shown in LNURL metadata. | 
