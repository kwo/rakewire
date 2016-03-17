package rest

import (
	"fmt"
	"github.com/gorilla/context"
	"log"
	"net/http"
	"rakewire/model"
	"rakewire/opml"
)

func (z *API) opmlExport(w http.ResponseWriter, req *http.Request) {

	user := context.Get(req, "user").(*model.User)

	var opmldoc *opml.OPML
	err := z.db.Select(func(tx model.Transaction) error {
		doc, err := opml.Export(tx, user)
		if err == nil {
			opmldoc = doc
		}
		return err
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set(hContentType, "application/xml")
	w.WriteHeader(http.StatusOK)
	err = opml.Format(opmldoc, w)
	if err != nil {
		log.Printf("%-7s %-7s Error formatting OPML: %s", logWarn, logName, err.Error())
	}

}

func (z *API) opmlImport(w http.ResponseWriter, req *http.Request) {

	user := context.Get(req, "user").(*model.User)

	opmldoc, err := opml.Parse(req.Body)
	if err != nil {
		message := fmt.Sprintf("Error parsing OPML: %s\n", err.Error())
		log.Printf("%-7s %-7s %s", logWarn, logName, message)
		w.Header().Set(hContentType, "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(message))
		return
	}

	err = z.db.Update(func(tx model.Transaction) error {
		return opml.Import(tx, user.ID, opmldoc)
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
