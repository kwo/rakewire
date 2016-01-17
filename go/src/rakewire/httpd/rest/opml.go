package rest

import (
	"github.com/gorilla/context"
	"log"
	"net/http"
	"rakewire/model"
	"rakewire/opml"
)

func (z *API) opmlExport(w http.ResponseWriter, req *http.Request) {

	user := context.Get(req, "user").(*model.User)

	doc, err := opml.Export(user, z.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set(hContentType, "application/xml")
	w.WriteHeader(http.StatusOK)
	err = opml.Format(doc, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (z *API) opmlImport(w http.ResponseWriter, req *http.Request) {

	err := func() error {
		user, err := z.db.UserGetByUsername("karl@ostendorf.com")
		if err != nil {
			return err
		}

		o, err := opml.Parse(req.Body)
		if err != nil {
			log.Printf("%-7s %-7s Error importing OPML: %s", logWarn, logName, err.Error())
			return err
		}

		replace := req.URL.Query().Get("replace") == "true"
		return opml.Import(user.ID, o, replace, z.db)

	}()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
