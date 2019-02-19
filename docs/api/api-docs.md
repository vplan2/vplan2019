# General API Documentation

---

## First Things First

### Parameters

Parameters with `default` formated names are **required** and parameters with *`italic`* formated names are ***optional***.

### Authentication

Generally, every endpoint requires an authorization.

The entry point for that will be the [`Authenticate`](#authenticate) endpoint.

For all other endpoints, there are 2 options to authenticate against the API:

**Authentication header**  
By sending an `Authentication` header with the raw token string as `value`, which you will get from the [`Authenticate`](#authenticate) endpoint. The `expire` value will say when the token validity will expire (as [UNIX timestamp](https://www.unixtimestamp.com/)).

**Session authentication**  
If you specify the authentication as `session` by parameter using the [`Authenticate`](#authenticate) endpoint, a cookie will be set which will authenticate you against the API while the cookie is valid.

### Status and Error Codes

The API uses the standart HTTP status codes, as defined in [RFC 2616](https://tools.ietf.org/html/rfc2616#section-10) and [RFC 6585](https://tools.ietf.org/html/rfc6585).

An error response from the API contains the status code as header and an error description as JSON body. Example error response:

```
< HTTP/1.1 400 Bad Request
< Content-Type: application/json
```
```json
{
  "error": {
    "code": 400,
    "message": "json: cannot unmarshal string into Go struct field authRequestData.session of type int"
  }
}
```

### Rate Limits

Rate limits are applied on a per-route basis, which means, that different API routs will count different rate limits. That also means, if you are curerntly rate-mimited on one endpoint, you can also use the other endpoints at this time in their specific rate limitations.

If you exeed a rate limit, you will get an error response as following:

```
< HTTP/1.1 429 Too many requests
< Content-Type: application/json
```
```json
{
  "error": {
    "code": 429,
    "message": "too many requests"
  }
}
```

---

## Endpoints

### Authenticate

> POST /api/authenticate/:USERNAME

#### Parameters

| Name | Type | Description |
|------|------|-------------|
| `password` | `string` | The password of the user |
| *`session`* | `int` | Specify if the login shoulb be treated as session creation which sets the authentication credentials as cookie. If this value is set `> 0`, you will not get an API key as response.<br/>`1` - basic session (valid for 1 hour)<br/>`2` - remembered session (valid for 30 days) |

#### Response

If `session` value is `0`

```
< HTTP/1.1 200 OK
< Content-Type: application/json
```
```json
{
    "token": "OTY2NDExMzE0MTU1MDU4NDA2OTMyOTEwMzkwMA==",
    "valid": 1553176069
}
```

If `session` value is larger than `0`

```
< HTTP/1.1 200 OK
< Content-Type: application/json
< Set-Cookie: main=MTU1MDU4NDEzNnxCQXdBQVRNPXwUbN6LLxTTxJ-eXp2SjGxnBg4o6E-IMFTUz1m2daa0aQ==; Path=/; Expires=Tue, 19 Feb 2019 14:48:56 GMT; Max-Age=3600
< Date: Tue, 19 Feb 2019 13:48:56 GMT
< Content-Length: 0
```