package webserver

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/vplan2/vplan2019/internal/database"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/vplan2/vplan2019/internal/auth"
	"github.com/vplan2/vplan2019/internal/logger"
)

/////////////////
// DATA STRUCT //
/////////////////

// authRequestData contains request data for
// POST /api/authenticate/:USERNAME
type authRequestData struct {
	Password string `json:"password"`
	Group    string `json:"group"`
	Session  int    `json:"session"`
}

// authResponseData contains the response data for
// POST /api/authenticate/:USERNAME
type authResponseData struct {
	Ident string      `json:"ident"`
	Ctx   interface{} `json:"ctx"`
}

// authTokenResposeData contains token string
// and expire time for token request response
type authTokenResposeData struct {
	*authResponseData
	Token  string    `json:"token"`
	Expire time.Time `json:"expire"`
}

////////////////////
// FRONTEND ROOTS //
////////////////////

func (s *Server) handlerFEMainRoot(w http.ResponseWriter, r *http.Request) {
	const indexPage = "/index.html"

	file := mux.Vars(r)["file"]
	if _, err := s.reqAuth.Authorize(r); err != nil {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {

		filepath := s.config.StaticFiles + "/" + file
		if !strings.HasSuffix(filepath, indexPage) {
			filepath += indexPage
		}

		if _, err = os.Stat(filepath); os.IsNotExist(err) {
			s.handlerErrorNotFound(w, r)
			return
		}

		http.ServeFile(w, r, s.config.StaticFiles+"/"+file)
	}
}

func (s *Server) handlerFELogin(w http.ResponseWriter, r *http.Request) {
	if _, err := s.reqAuth.Authorize(r); err != nil {
		http.ServeFile(w, r, s.config.StaticFiles+"/login/index.html")
	} else {
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

/////////
// API //
/////////

// POST /api/authenticate/:USERNAME
func (s *Server) handlerAPIAuthenticate(w http.ResponseWriter, r *http.Request) {
	if !s.limiter.Check("authenticate", w, r) {
		return
	}

	urlParams := mux.Vars(r)
	uname := urlParams["username"]

	ipaddress := r.RemoteAddr
	useragent := r.Header.Get("User-Agent")
	if useragent == "" {
		s.handlerAPIBadRequestError(w, r, "User-Agent header must be set")
		return
	}

	reqData := new(authRequestData)
	if err := s.parseJSONBody(r.Body, reqData); err != nil {
		s.handlerAPIBadRequestError(w, r, err.Error())
		return
	}

	passwd := reqData.Password
	if passwd == "" {
		s.handlerAPIUnauthorizedError(w, r)
		return
	}

	authData, err := s.authProvider.Authenticate(uname, reqData.Group, passwd)
	if err != nil {
		s.handlerAPIUnauthorizedError(w, r)
		return
	}

	// Just to ensure we do not run into an runtime error
	// later on using this object
	if authData == nil {
		authData = new(auth.Response)
	}

	respData := &authResponseData{
		Ident: authData.Ident,
		Ctx:   authData.Ctx,
	}

	if reqData.Session > 0 {
		var session *sessions.Session
		session, _ = s.store.Get(r, auth.MainSessionName)
		session.Values["ident"] = authData.Ident

		if uname == s.config.TVUser {
			session.Options.MaxAge = 157680000 // == 5 years
		} else if reqData.Session > 1 {
			session.Options.MaxAge = s.config.Sessions.RememberMaxAge
		}

		err := session.Save(r, w)
		if err != nil {
			s.handlerAPIInternalError(w, r, err)
			return
		}
	} else {
		token, expire, err := s.tokenManager.Set(authData.Ident)
		if err != nil {
			s.handlerAPIInternalError(w, r, err)
		} else {
			s.db.InsertLogin(database.LoginTypeToken, authData.Ident, useragent, ipaddress)
			jsonResponse(w, http.StatusOK, authTokenResposeData{
				Token:            token,
				Expire:           expire,
				authResponseData: respData,
			})
		}
		return
	}

	s.db.InsertLogin(database.LoginTypeWebInterface, authData.Ident, useragent, ipaddress)
	jsonResponse(w, http.StatusOK, respData)
}

// POST /api/logout
func (s *Server) handlerAPILogout(w http.ResponseWriter, r *http.Request) {
	if !s.limiter.Check("logout", w, r) {
		return
	}

	w.Header().Set("Set-Cookie", auth.MainSessionName+"=deleted; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT")
	jsonResponse(w, http.StatusOK, nil)
}

// GET /api/logins
func (s *Server) handlerAPIGetLogins(w http.ResponseWriter, r *http.Request) {
	if !s.limiter.Check("getLogins", w, r) {
		return
	}

	ident := s.reqAuth.Check(w, r)
	if ident == "" {
		return
	}

	urlQuery := r.URL.Query()

	_t := urlQuery.Get("time")
	t, ok := parseTimeRFC3339(w, _t, false)
	if !ok {
		return
	}

	limit := 20
	_limit := urlQuery.Get("limit")
	if _limit != "" {
		var err error
		limit, err = strconv.Atoi(_limit)
		if err != nil {
			s.handlerAPIBadRequestError(w, r, "limit must be a valid number")
			return
		}
	}

	logins, err := s.db.GetLogins(ident, t, limit)
	if err != nil {
		s.handlerAPIInternalError(w, r, err)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"data": logins,
	})
}

// GET /api/vplan
func (s *Server) handlerAPIGetVPlan(w http.ResponseWriter, r *http.Request) {
	if !s.limiter.Check("getVPlan", w, r) {
		return
	}

	ident := s.reqAuth.Check(w, r)
	if ident == "" {
		return
	}

	var err error

	reqQuery := r.URL.Query()

	class := reqQuery.Get("class")
	_ignoreSettings := reqQuery.Get("ignoreSettings")
	ignoreSettings := 0
	if _ignoreSettings != "" {
		ignoreSettings, err = strconv.Atoi(_ignoreSettings)
		if err != nil {
			s.handlerAPIBadRequestError(w, r, "invalid format for ignoreSettings: must be a number")
			return
		}
	}

	_t := reqQuery.Get("time")

	var t time.Time
	var ok bool
	if _t == "" {
		t = time.Now()
	} else {
		if t, ok = parseTimeRFC3339(w, _t, true); !ok {
			return
		}
	}

	if class == "" && ignoreSettings <= 0 {
		uSettings, ok, err := s.db.GetUserSettings(ident)
		if err != nil {
			s.handlerAPIInternalError(w, r, err)
			return
		}
		if ok {
			class = uSettings.Class
		}
	}

	vplans, err := s.db.GetVPlans(class, t)
	if err != nil {
		s.handlerAPIInternalError(w, r, err)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"data": vplans,
	})
}

// GET /api/newsticker
func (s *Server) handlerAPIGetNewsTicker(w http.ResponseWriter, r *http.Request) {
	if !s.limiter.Check("getNewsTicker", w, r) {
		return
	}

	ident := s.reqAuth.Check(w, r)
	if ident == "" {
		return
	}

	_t := r.URL.Query().Get("time")

	var t time.Time
	var ok bool
	if _t == "" {
		t = time.Now().AddDate(0, -1, 0)
	} else {
		if t, ok = parseTimeRFC3339(w, _t, true); !ok {
			return
		}
	}

	entries, err := s.db.GetNewsTicker(t)
	if err != nil {
		s.handlerAPIInternalError(w, r, err)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"data": entries,
	})
}

// GET /api/settings
func (s *Server) handleAPIGetUserSettings(w http.ResponseWriter, r *http.Request) {
	if !s.limiter.Check("getUserSettings", w, r) {
		return
	}

	ident := s.reqAuth.Check(w, r)
	if ident == "" {
		return
	}

	settings, _, err := s.db.GetUserSettings(ident)
	if err != nil {
		s.handlerAPIInternalError(w, r, err)
		return
	}

	jsonResponse(w, http.StatusOK, settings)
}

// POST /api/settings
func (s *Server) handleAPISetUserSettings(w http.ResponseWriter, r *http.Request) {
	if !s.limiter.Check("setUserSettings", w, r) {
		return
	}

	ident := s.reqAuth.Check(w, r)
	if ident == "" {
		return
	}

	updateSettings := new(database.UserSetting)
	err := s.parseJSONBody(r.Body, updateSettings)
	if err != nil {
		s.handlerAPIBadRequestError(w, r, err.Error())
		return
	}

	err = s.db.SetUserSetting(ident, updateSettings)
	if err != nil {
		s.handlerAPIInternalError(w, r, err)
		return
	}

	jsonResponse(w, http.StatusOK, nil)
}

// POST /api/test
// Just for testing purposes
func (s *Server) handlerAPITest(w http.ResponseWriter, r *http.Request) {
	if !s.limiter.Check("test", w, r) {
		return
	}

	ident := s.reqAuth.Check(w, r)
	if ident == "" {
		return
	}

	logger.Debug("auth test: %s", ident)
}

//////////////////
// HELPER FUNCS //
//////////////////

// parseTimeRFC3339 tries to parse a string to a time.Time struct by using the RFC3339
// format. If the parsing fails, a bad request response will be written to the
// response writer and the function returns a '0' time with false as second return value.
//   w:              HTTP response writer
//   rawString:      time stamp to parse
//   mustNotBeEmpty: if this is set to 'true' and rawString is '""', a bad request will be
//                   send as response. Else, the function will return a '0' time struct
//                   with true as second return value.
func parseTimeRFC3339(w http.ResponseWriter, rawString string, mustNotBeEmpty bool) (time.Time, bool) {
	if rawString == "" && !mustNotBeEmpty {
		return time.Time{}, true
	}

	t, err := time.Parse(time.RFC3339, rawString)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest,
			apiError(http.StatusBadRequest, "time format is not RFC3339"))
		return time.Time{}, false
	}

	return t, true
}

////////////////////
// ERROR HANDLERS //
////////////////////

func (s *Server) handlerErrorNotFound(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, s.config.StaticFiles+"/404.html")
}

func (s *Server) handlerAPIInternalError(w http.ResponseWriter, r *http.Request, err error) {
	jsonResponse(w, http.StatusInternalServerError, apiError(http.StatusInternalServerError, err.Error()))
}

func (s *Server) handlerAPIUnauthorizedError(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusUnauthorized, apiError(http.StatusUnauthorized, ""))
}

func (s *Server) handlerAPIRateLimitError(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusTooManyRequests, apiError(http.StatusTooManyRequests, ""))
}

func (s *Server) handlerAPIBadRequestError(w http.ResponseWriter, r *http.Request, msg string) {
	jsonResponse(w, http.StatusBadRequest, apiError(http.StatusBadRequest, msg))
}
