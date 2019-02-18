package webserver

import (
	"fmt"
	"html/template"
	"net/http"

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
	Session  int    `json:"session"`
}

//////////////
// FRONTEND //
//////////////

func (s *Server) handlerMainPage(w http.ResponseWriter, r *http.Request) {
	t := template.New("index.html")
	_, err := t.ParseFiles("./web/views/index.html")
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

	authData, err := s.authProvider.Authenticate(uname, passwd)
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, apiError(http.StatusUnauthorized, ""))
		return
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

	}

	jsonResponse(w, http.StatusOK, nil)
}

func (s *Server) handlerAPITest(w http.ResponseWriter, r *http.Request) {
	if !s.limiter.Check("test", w, r) {
		return
	}

	session, err := s.store.Get(r, "main")
	fmt.Println(err)
	fmt.Println(session.Values)
}

// ERROR HANDLERS

func (s *Server) handlerAPIRateLimitError(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusTooManyRequests, apiError(http.StatusTooManyRequests, ""))
}
