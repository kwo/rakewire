package rest

import (
	"log"
	"net/http"
	"rakewire/opml"
)

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

		return opml.Import(user.ID, o, true, z.db)

	}()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
