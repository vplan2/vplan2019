# General API Documentation

---

## First Things First

#### Parameters

Parameters with a `default` formated names are **required** and parameters with *`italic`* formated names are ***optional***.

#### Authentication

Generally, every endpoint except the `Authenticate` endpoint requires an authorization.

There are 2 options to authenticate against the API:

**Authentication header**  
By sending an `Authentication` header with the raw token string as `value`, which you will get from the `Authenticate` endpoint. The `valid` value will say for how long the token will be valid (in seconds).

**Session authentication**  
If you specify the authentication as `session` by parameter using the `Authenticate` endpoint, a cookie will be set which will authenticate you against the API while the cookie is valid.

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

```
Status: 200 OK
```
```json
// Only if 'session' value was 0
{
    "token": "UEtRUmg3djdibGtmOEEweEJvdWV6ZU82NFhyOXNXWk1TVGtKMzAzdmhPbDNRU2ljUU9EaWVhZVpYazROWHR6RQo=",
    "valid": 3600
}
```