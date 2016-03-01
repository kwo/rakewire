package rest

import (
	"fmt"
	"github.com/gorilla/context"
	"log"
	"net/http"
	"rakewire/model"
)

func (z *API) opmlExport(w http.ResponseWriter, req *http.Request) {

	user := context.Get(req, "user").(*model.User)

	var opml *model.OPML
	err := z.db.Select(func(tx model.Transaction) error {
		doc, err := model.OPMLExport(user, tx)
		if err == nil {
			opml = doc
		}
		return err
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set(hContentType, "application/xml")
	w.WriteHeader(http.StatusOK)
	err = model.OPMLFormat(opml, w)
	if err != nil {
		log.Printf("%-7s %-7s Error formatting OPML: %s", logWarn, logName, err.Error())
	}

}

func (z *API) opmlImport(w http.ResponseWriter, req *http.Request) {

	user := context.Get(req, "user").(*model.User)

	opml, err := model.OPMLParse(req.Body)
	if err != nil {
		message := fmt.Sprintf("Error parsing OPML: %s\n", err.Error())
		log.Printf("%-7s %-7s %s", logWarn, logName, message)
		w.Header().Set(hContentType, "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(message))
		return
	}

	replace := req.URL.Query().Get("replace") == "true"
	err = z.db.Update(func(tx model.Transaction) error {
		return model.OPMLImport(user.ID, opml, replace, tx)
	})

	if err != nil {
		message := fmt.Sprintf("Error importing OPML: %s\n", err.Error())
		log.Printf("%-7s %-7s %s", logWarn, logName, message)
		w.Header().Set(hContentType, "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(message))
		return
	}

	w.Header().Set(hContentType, "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))

}
