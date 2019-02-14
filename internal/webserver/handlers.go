package webserver

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/zekroTJA/vplan2019/internal/logger"
)

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
	urlParams := mux.Vars(r)
	_, ok := urlParams["username"]
	if !ok {
		jsonResponse(w, http.StatusBadRequest, nil)
		return
	}
}
