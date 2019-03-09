# General API Documentation

---

## First Things First

### Parameters

Parameters with `default` formated names are **required** and parameters with *`italic`* formated names are ***optional***.

### Authentication

Generally, every endpoint requires an authorization.

The entry point for that will be the [`Authenticate`](#authenticate) endpoint.

For all other endpoints, there are 2 options to authenticate against the API:

> TODO: I) anpassen wie II)

**I) Session authentication**  
If you specify the authentication as `session` by parameter using the [`Authenticate`](#authenticate) endpoint, a cookie will be set which will authenticate you against the API while the cookie is valid. Make sure, that you application is capable to save cookies and send them like described in [RFC 6265, section 5.4](https://tools.ietf.org/html/rfc6265#section-5.4) if you are using this method.

**I) Authentication header**  
Otherwise, no session cookie will be set. Therefore, an API token will be generated and returned as `token` value in the response body in combination with an expire [timestamp](#data-formats) defined after the `expire` key. In order to authenticate against the API, you must send an `Authorization` header with the token string as value on every following request.

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

Rate limits are applied on a per-route basis, which means that different API routs will count different rate limits. That also means, if you are curerntly rate-limited on one endpoint, you can also use the other endpoints at this time in their specific rate limitations.

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

### Data Formats

**Timestamps**  
All time formats of request data and response data must be formatted and interpreted as [RFC 3339](https://tools.ietf.org/html/rfc3339) formatted timestamp.

---

## Endpoints

- [Authenticate](#authenticate)  
  `POST /api/authenticate/:USERNAME`
- [Logout](#logout)  
  `POST /api/logout`
- [Get Logins](#get-logins)  
  `GET /api/logins`

- [Get VPlans](#get-vplans)  
  `GET /api/vplan`

- [Get User Settings](#get-user-settings)  
  `GET /api/settings`
- [Set User Settings](#set-user-settings)  
  `POST /api/settings`

---

### Authenticate

> POST /api/authenticate/:USERNAME

#### Parameters

| Name | Type | Description |
|------|------|-------------|
| `password` | `string` | The password of the user |
| *`group`* | `string` | If the server requires a group to authenticate, you can specify this here |
| *`session`* | `int` | Specify if the login should be treated as session creation which sets the authentication credentials as cookie. If this value is set `> 0`, you will not get an API key as response.<br/>`1` - basic session (valid for 1 hour)<br/>`2` - remembered session (valid for 30 days) |

#### Response

If `session` value is `0`

```
< HTTP/1.1 200 OK
< Content-Type: application/json
```
```json
{
  "ident": "cn=mustermax,dc=example,dc=de",
  "ctx": {
    "cn": [
      "mustermax"
    ],
    "ou": [
      "user"
    ]
  },
  "token": "Nzc2Nzk2Mzc4MDE1NTE3Nzk4MDk5NjA0OTQ2MDA=",
  "expire": "2019-04-05T15:09:57.9536976+02:00"
}
```

If `session` value is larger than `0`

```
< HTTP/1.1 200 OK
< Content-Type: application/json
< Set-Cookie: main=MTU1MDU4NDEzNnxCQXdBQVRNPXwUbN6LLxTTxJ-eXp2SjGxnBg4o6E-IMFTUz1m2daa0aQ==; Path=/; Expires=Tue, 19 Feb 2019 14:48:56 GMT; Max-Age=3600
```

```json
{
  "ident": "cn=mustermax,ou=user,dc=zekro,dc=de",
  "ctx": {
    "cn": [
      "mustermax"
    ],
    "ou": [
      "user"
    ]
  }
}
```

### Logout

> POST /api/logout

#### Parameters

*No parameters required.*

#### Response

> This endpoint does not respond with any body information. Also, it does not check if you already have set a session cookie. It will just set the session cookie as deleted and expired, so the cookie will be deleted on session ending.

```
< HTTP/1.1 200 OK
```

### Get Logins

> GET /api/logins

#### Parameters
> Parameters must be passed by *(URL encoded)* URL parameters.

| Name | Type | Description |
|------|------|-------------|
| *`time`* | `string` | [RFC 3339](https://tools.ietf.org/html/rfc3339) encoded timestamp after which logins are requested |

#### Response

Response contains a `type` value of the logins which must be interpreted as following:

| Type | Description |
|------|-------------|
| `0` | **Web Interface Login** - A default login over the web interface mask using session cookies to authenticate |
| `1` | **Token Creation** - An API token generated by the authenticate endpoint |

```
< HTTP/1.1 200 OK
< Content-Type: application/json
```
```json
{
  "data": [
    {
      "ident": "cn=mustermax,dc=example,dc=de",
      "timestamp": "2019-03-07T08:40:44Z",
      "type": 0,
      "useragent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:65.0) Gecko/20100101 Firefox/65.0",
      "ipaddress": "123.45.67.89:53321"
    },
    {
      "ident": "cn=mustermax,dc=example,dc=de",
      "timestamp": "2019-03-07T10:04:16Z",
      "type": 0,
      "useragent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36",
      "ipaddress": "13.37.69.111:54168"
    },
    {
      "ident": "cn=mustermax,dc=example,dc=de",
      "timestamp": "2019-03-07T11:14:35Z",
      "type": 1,
      "useragent": "MyCustomWebApplicationUsingTokenSystem",
      "ipaddress": "13.37.69.111:54168"
    }
  ]
}
```

### Get VPlans

> GET /api/vplan

#### Parameters
> Parameters must be passed by *(URL encoded)* URL parameters.

| Name | Type | Description |
|------|------|-------------|
| *`time`* | `string` | [RFC 3339](https://tools.ietf.org/html/rfc3339) encoded timestamp after which VPlans are requested |
| *`class`* | `string` | Name of the class of which the VPlan entries will be filtered |

#### Response

```
< HTTP/1.1 200 OK
< Content-Type: application/json
```
```json
{
  "data": [
    {
      "id": 1397,
      "date_edit": "2019-01-18T09:17:47Z",
      "date_for": "2019-03-08T00:00:00Z",
      "block": "X",
      "header": "",
      "footer": "",
      "entries": [
        {
          "id": 26352,
          "vplan_id": 1397,
          "class": "UL17X",
          "time": "7./8.",
          "messures": "LF8 in 034 vom 17.01.",
          "responsible": "Hr. A"
        }
      ]
    },
    {
      "id": 1408,
      "date_edit": "2019-02-15T10:09:12Z",
      "date_for": "2019-03-04T00:00:00Z",
      "block": "F",
      "header": "Montag, den 04.03.2019\r\nabwesend: Hr. W",
      "footer": "IO16F Projekt",
      "entries": [
        {
          "id": 26776,
          "vplan_id": 1408,
          "class": "TZ17F",
          "time": "1./2. ",
          "messures": "Gr. W Ausfall ETCS",
          "responsible": "Hr. W"
        },
        {
          "id": 26778,
          "vplan_id": 1408,
          "class": "WQ16F",
          "time": "3./4.",
          "messures": "ganze Klasse RDSE in H1",
          "responsible": "Hr. J"
        }
      ]
    },
    {
      "id": 1410,
      "date_edit": "2019-02-26T10:47:13Z",
      "date_for": "2019-03-05T00:00:00Z",
      "block": "B",
      "header": "Dienstag, den 05.03. 2019\r\nabwesend: ",
      "footer": "",
      "entries": null
    }
  ]
}
```

### Get User Settings

> GET /api/settings

#### Parameters
> Parameters must be passed by *(URL encoded)* URL parameters.

#### Response

> Unset values are default type values, e.g. `""` for strings or `0` for integers.

```
< HTTP/1.1 200 OK
< Content-Type: application/json
```
```json
{
  "ident": "cn=mustermax,dc=example,dc=de",
  "class": "ITF17C",
  "theme": "dark",
  "edited": "2019-03-07T16:22:19Z"
}
```

### Set User Settings

> POST /api/settings

#### Parameters

> Only passed arguments will be updated in the user settings. If you want to set the value to the default initialization value of the value type, use the defined reset value.

| Name | Type | Description |
|------|------|-------------|
| *`class`* | `string` | Users class name (`reset` to reset value) |
| *`theme`* | `string` | UI theme (`reset` to reset value) |

#### Response

```
< HTTP/1.1 200 OK
```

---

© 2019 Justin Trommler, Richard Heidenreich, Ringo Hoffmann  
Covered by MIT Licence.

The API structure and design is inspired by the REST API's of [discordapp.com](https://discordapp.com/developers/docs/intro) and [github.com](https://developer.github.com/v3/).