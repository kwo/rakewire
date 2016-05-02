package httpd

import (
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"net/http"
	"rakewire/auth"
	"rakewire/model"
	"rakewire/opml"
)

type oddballs struct {
	db model.Database
}

func (z *oddballs) router() *mux.Router {

	router := mux.NewRouter()

	router.Path("/x/subscriptions.opml").Methods(mGet).HandlerFunc(z.opmlExport)
	router.Path("/x/subscriptions.opml").Methods(mPut).HandlerFunc(z.opmlImport)

	return router

}

func (z *oddballs) opmlExport(w http.ResponseWriter, req *http.Request) {

	user := context.Get(req, "user").(*auth.User)

	var opmldoc *opml.OPML
	err := z.db.Select(func(tx model.Transaction) error {
		if u := model.U.GetByUsername(tx, user.Name); u != nil {
			doc, errExport := opml.Export(tx, u)
			if errExport != nil {
				return errExport
			}
			opmldoc = doc
			return nil
		}
		return fmt.Errorf("User not found: %s", user.Name)
	})

	if err != nil {
		message := fmt.Sprintf("Error creating OPML: %s\n", err.Error())
		log.Debugf("%s", message)
		w.Header().Set(hContentType, "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(message))
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

	user := context.Get(req, "user").(*auth.User)

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
		if u := model.U.GetByUsername(tx, user.Name); u != nil {
			return opml.Import(tx, u.ID, opmldoc)
		}
		return fmt.Errorf("User not found: %s", user.Name)
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
