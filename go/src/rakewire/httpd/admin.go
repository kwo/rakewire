package httpd

import (
	"log"
	"net/http"
)

func (z *Service) repairDatabase(w http.ResponseWriter, req *http.Request) {

	err := z.Database.Repair()
	if err != nil {
		log.Printf("%s %s Error in db.Repair: %s\n", logWarn, logName, err.Error())
		http.Error(w, "Cannot repair database.", http.StatusInternalServerError)
		return
	}

	sendOK(w, req)

}
