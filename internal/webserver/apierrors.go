package webserver

var stdErrMsgs = map[int]string{
	400: "bad request",
	401: "unauthorized",
	403: "forbidden",
	429: "too many requests",
}

type apiErrorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func apiError(code int, msg string) map[string]interface{} {
	if msg == "" {
		msg = stdErrMsgs[code]
	}

	return map[string]interface{}{
		"error": apiErrorBody{
			Code:    code,
			Message: msg,
		},
	}
}
