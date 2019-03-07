package webserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zekroTJA/vplan2019/internal/database"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/zekroTJA/vplan2019/internal/auth"
	"github.com/zekroTJA/vplan2019/internal/logger"
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
	file := mux.Vars(r)["file"]
	if _, err := s.reqAuth.Authorize(r); err != nil {
		http.ServeFile(w, r, s.config.StaticFiles+"/login.html")
	} else {
		fmt.Println(file)
		http.ServeFile(w, r, s.config.StaticFiles+"/"+file)
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
		jsonResponse(w, http.StatusBadRequest, apiError(http.StatusBadRequest, "User-Agent header must be set"))
		return
	}

	reqData := new(authRequestData)
	if err := s.parseJSONBody(r.Body, reqData); err != nil {
		jsonResponse(w, http.StatusBadRequest, apiError(http.StatusBadRequest, err.Error()))
		return
	}

	passwd := reqData.Password
	if passwd == "" {
		jsonResponse(w, http.StatusBadRequest, apiError(http.StatusBadRequest, ""))
		return
	}

	authData, err := s.authProvider.Authenticate(uname, reqData.Group, passwd)
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, apiError(http.StatusUnauthorized, ""))
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
		if reqData.Session > 1 {
			session.Options.MaxAge = s.config.Sessions.RememberMaxAge
		}
		err := session.Save(r, w)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, apiError(http.StatusInternalServerError, err.Error()))
			return
		}
	} else {
		token, expire, err := s.tokenManager.Set(authData.Ident)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, apiError(http.StatusInternalServerError, err.Error()))
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

	_t := r.URL.Query().Get("time")
	t, ok := parseTimeRFC3339(w, _t, false)
	if !ok {
		return
	}

	logins, err := s.db.GetLogins(ident, t)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, apiError(http.StatusInternalServerError, err.Error()))
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

	reqQuery := r.URL.Query()

	class := reqQuery.Get("class")
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

	vplans, err := s.db.GetVPlans(class, t)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError,
			apiError(http.StatusInternalServerError, err.Error()))
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"data": vplans,
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
		jsonResponse(w, http.StatusInternalServerError,
			apiError(http.StatusInternalServerError, err.Error()))
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
		jsonResponse(w, http.StatusBadRequest,
			apiError(http.StatusBadRequest, err.Error()))
		return
	}

	err = s.db.SetUserSetting(ident, updateSettings)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError,
			apiError(http.StatusInternalServerError, err.Error()))
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

func (s *Server) handlerAPIInternalError(w http.ResponseWriter, r *http.Request, err error) {
	jsonResponse(w, http.StatusInternalServerError, apiError(http.StatusInternalServerError, err.Error()))
}

func (s *Server) handlerAPIUnauthorizedError(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusUnauthorized, apiError(http.StatusUnauthorized, ""))
}

func (s *Server) handlerAPIRateLimitError(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusTooManyRequests, apiError(http.StatusTooManyRequests, ""))
}
