package httpd

import (
	"net/http"
)

func (z *Service) repairDatabase(w http.ResponseWriter, req *http.Request) {

	err := z.Database.Repair()
	if err != nil {
		logger.Warnf("Error in db.Repair: %s\n", err.Error())
		http.Error(w, "Cannot repair database.", http.StatusInternalServerError)
		return
	}

	sendOK(w, req)

}
