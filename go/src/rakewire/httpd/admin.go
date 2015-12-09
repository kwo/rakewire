package httpd

import (
	"log"
	"net/http"
)

func (z *Service) repairDatabase(w http.ResponseWriter, req *http.Request) {

	err := z.database.Repair()
	if err != nil {
		log.Printf("%-7s %-7s Error in db.Repair: %s", logWarn, logName, err.Error())
		http.Error(w, "Cannot repair database.", http.StatusInternalServerError)
		return
	}

	sendOK(w, req)

}
