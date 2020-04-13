package handlers

import (
	"net/http"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/pkg/srverror"
	"git.maxset.io/web/knaxim/pkg/srvjson"
	"github.com/gorilla/mux"
)

// AttachAcronym is to add api paths related to acronyms
func AttachAcronym(r *mux.Router) {
	r = r.NewRoute().Subrouter()
	r.Use(srvjson.JSONResponse)
	r.Use(ConnectDatabase)
	r.Use(UserCookie)

	r.HandleFunc("/{acronym}", getAcronym).Methods("GET")
}

func getAcronym(out http.ResponseWriter, r *http.Request) {
	w, ok := out.(*srvjson.ResponseWriter)
	if !ok {
		panic(srverror.Basic(500, "Server Error", "expecting *srvjson.ResponseWriter"))
	}
	vals := mux.Vars(r)
	matches, err := r.Context().Value(types.ACRONYM).(database.Acronymbase).Get(vals["acronym"])
	if err != nil {
		panic(err)
	}
	w.Set("matched", matches)
}
