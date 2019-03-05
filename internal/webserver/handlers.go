package webserver

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/zekroTJA/vplan2019/internal/auth"

	"github.com/gorilla/mux"

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

type authResponseData struct {
	Ident string      `json:"ident"`
	Ctx   interface{} `json:"ctx"`
}

// authTokenResposeData contains token string
// and expire time for token request response
type authTokenResposeData struct {
	*authResponseData
	Token  string `json:"token"`
	Expire int64  `json:"expire"`
}

//////////////
// FRONTEND //
//////////////

// Handler for root page
func (s *Server) handlerMainPage(w http.ResponseWriter, r *http.Request) {
	t := template.New("index.html")
	_, err := t.ParseFiles(s.config.StaticFiles + "/web/views/index.html")
	if err != nil {
		logger.Error("failed parsing HTML template: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, struct {
		Test string
	}{
		"testData",
	})
	if err != nil {
		logger.Error("failed parsing HTML template: ", err)
		return
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
	uname, ok := urlParams["username"]
	if !ok {
		jsonResponse(w, http.StatusBadRequest, apiError(http.StatusBadRequest, ""))
		return
	}

	data := new(authRequestData)
	if err := s.parseJSONBody(r.Body, data); err != nil {
		jsonResponse(w, http.StatusBadRequest, apiError(http.StatusBadRequest, err.Error()))
		return
	}

	passwd := data.Password
	if passwd == "" {
		jsonResponse(w, http.StatusBadRequest, apiError(http.StatusBadRequest, ""))
		return
	}

	authData, err := s.authProvider.Authenticate(uname, data.Group, passwd)
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

	if data.Session > 0 {
		session, err := s.store.Get(r, "main")
		session.Values["ident"] = authData.Ident
		if data.Session > 1 {
			session.Options.MaxAge = s.config.Sessions.RememberMaxAge
		}
		session.Save(r, w)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, apiError(http.StatusInternalServerError, err.Error()))
			return
		}
	} else {
		token, expire, err := s.tokenManager.Set(authData.Ident)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, apiError(http.StatusInternalServerError, err.Error()))
		} else {
			jsonResponse(w, http.StatusOK, authTokenResposeData{
				Token:            token,
				Expire:           expire.Unix(),
				authResponseData: respData,
			})
		}
		return
	}

	jsonResponse(w, http.StatusOK, respData)
}

// POST /api/test
// Just for testing purposes
func (s *Server) handlerAPITest(w http.ResponseWriter, r *http.Request) {
	if !s.limiter.Check("test", w, r) {
		return
	}
	// session, err := s.store.Get(r, "main")
	// fmt.Println(err)
	// fmt.Println(session.Values)

	// fmt.Println(s.db.SetUserAPIToken("testUser", "testToken", time.Now().Add(10*time.Minute)))
	// fmt.Println(s.db.DeleteUserAPIToken("testUser"))

	fmt.Println(s.reqAuth.Check(w, r))
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
