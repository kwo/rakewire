package httpd

import (
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"net/http"
	"rakewire/model"
	"rakewire/opml"
)

type oddballs struct {
	db model.Database
}

func (z *oddballs) router() *mux.Router {

	router := mux.NewRouter()

	router.Path("/subscriptions.opml").Methods(mGet).HandlerFunc(z.opmlExport)
	router.Path("/subscriptions.opml").Methods(mPut).HandlerFunc(z.opmlImport)

	return router

}

func (z *oddballs) opmlExport(w http.ResponseWriter, req *http.Request) {

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
		log.Debugf("Error formatting OPML: %s", err.Error())
	}

}

func (z *oddballs) opmlImport(w http.ResponseWriter, req *http.Request) {

	user := context.Get(req, "user").(*model.User)

	opmldoc, err := opml.Parse(req.Body)
	if err != nil {
		message := fmt.Sprintf("Error parsing OPML: %s\n", err.Error())
		log.Debugf("%s", message)
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
		log.Debugf("%s", message)
		w.Header().Set(hContentType, "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(message))
		return
	}

	w.Header().Set(hContentType, "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))

}
