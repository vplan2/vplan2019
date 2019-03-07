+++
title = "Backend"
weight = 3
sort_by = "weight"
insert_anchor_links = "right"
+++

## Software:
- Development:
  - [Notepad++] / [Geany]
  - [Visual Studio Code] with Go-Package
  - [Insomia] a [REST] client
- Production:
  - Server or Hostsystem (Windows / Linux)

## Server / CGI:
- (required) VPLAN2 [REST]-API
  - written in: [Go]
  - (optional) with [go-daemon]
- (optional) Web Push Server
  - written in: [Rust] (cargo tools: fmt, clippy, rls, update?)
  - (optional) with [rust-daemon]

## Tasks:
- Authentification
- Autorisation
- [TLS] as spezified in [RFC 5246] or [RFC 8446]
- Application Programming Interface over 
  - [REST]
  - [CLI]? (like mysql or [sysinfo] example)
- HTTP Status Codes
  - [426] Browser support etc.
  - [418] Wrong Usage (specified in [RFC 2324])

## Interface / REST:

```bash
$ curl -X GET <host>/api/vplan
<json-output with vplan infos>
$ curl -X GET <host>/rss/releases
<rss-output with release notes>
```

```bash
$ vplan -u <user> -p <pass>
Welcome to Vplan 2.0 CLI
> /help
<help-output>
> /GET /api/vplan
<json-output>
```
