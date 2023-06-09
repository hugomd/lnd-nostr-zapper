# lnd-nostr-zapper

Self-host your own [zap](https://nostr.how/en/zaps) server to receive zaps
directly to your own [Umbrel](https://umbrel.com)/LND node!

Heavily inspired by [fiatjaf/bridgeaddr](https://github.com/fiatjaf/bridgeaddr/tree/master) and [jb55/cln-nostr-zapper](https://github.com/jb55/cln-nostr-zapper/tree/master).

# Features
Implements [LUD-06](https://github.com/lnurl/luds/blob/luds/06.md) (LNURLp) and
[NIP-57](https://github.com/nostr-protocol/nips/blob/master/57.md).

* [x] Connect to LND via Tor
* [x] Publish zap receipts
* [x] Configure default relays for zap receipts
* [x] Easy deployment with Fly.io
* [x] Public Docker image

# Installation

A Docker container image is available via [GitHub Container Registry](https://github.com/hugomd/lnd-nostr-zapper/pkgs/container/lnd-nostr-zapper):
```
docker pull ghcr.io/hugomd/lnd-nostr-zapper:e8de5840
```

# Deploying

## Fly.io

[Fly](https://fly.io) makes it easy to run Docker containers, and has a generous 
free tier.

Install `flyctl`:

```bash
brew install flyctl
```

Clone this repository:
```bash
git clone git@github.com:hugomd/lnd-nostr-zapper.git && cd lnd-nostr-zapper
```

### Set environment variables

Update the environment variables in [`fly.toml`](./fly.toml) (e.g. `LND_HOST`).

Most of the environment variables can be found on your Umbrel node by SSH'ing 
in and running the commands below:

To find your Tor hostname, run:
```bash
cat /home/umbrel/umbrel/tor/data/app-lightning-rest/hostname
```

`LND_HOST` should be of the form:
```
https://TOR_HOSTNAME:8080
```

To bake a macaroon run:
```bash
docker exec lightning_lnd_1 lncli bakemacaroon invoices:read invoices:write
```

### Deploy

Launch the application:
```bash
flyctl launch
```

### Set secrets

Set secret values, which will be available to the container as runtime
environment variables:
```bash
flyctl secrets set NOSTR_KEY="NOSTR_PRIVATE_KEY_HERE"
flyctl secrets set LND_MACAROON="LND_MACAROON_HERE"
```

This should deploy lnd-nostr-zapper to a `fly.dev` domain, which you can use to 
receive zaps! ⚡️

### Test!

You can see it in action by running:
```bash
flyctl open "/.well-known/lnurlp/capybara"
```

# Configuration
Configuration is done via environment variables:

| Environment Variable  | Required | Default value | Description |
| --------------------- | -------- | ------------- | ----------- |
| `HOST`                | false    | `0.0.0.0`      | The host to bind the HTTP server to. | 
| `PORT`                | false    | `8080`         | The port to bind the HTTP server to. | 
| `DOMAIN`              | true     | N/A            | The domain associated with the server, used in the LNURL callback. | 
| `LND_HOST`            | true     | N/A            | URL pointing to LND, followed by the REST API port. E.g. `https://example.onion:8080` | 
| `LND_MACAROON`        | true     | N/A            | An invoice read/write macaroon for auth with LND. | 
| `LND_CERT`            | false    | `""`           | Optional self-signed certificate to call LND. Don't set this if you run Umbrel. | 
| `NOSTR_KEY`           | true     | N/A            | Nostr private key, used to publish zap receipts. | 
| `COMMENT_LENGTH`      | false    | `0`            | Maximum length of associated comments, sent via webhook. | 
| `WEBHOOK_URL`         | false    | `""`           | URL to call after successful payment. | 
| `DESCRIPTION`         | false    | `"Send sats!"` | Description shown in LNURL metadata. | 
| `IMAGE_URL`           | false    | `""`           | Optional avatar URL shown in LNURL metadata. | 
| `RELAYS`              | false    | `""`           | Optional comma separated list of relays to publish zaps to, of the form: `wss://relay.damus.io,wss://brb.io`. | 
