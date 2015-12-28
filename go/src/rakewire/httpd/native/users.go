package native

import (
	"encoding/json"
	"log"
	"net/http"
)

func (z *API) usersGet(w http.ResponseWriter, req *http.Request) {

	rsp := &User{
		ID:       0,
		Username: "username",
	}

	w.Header().Set(hContentType, mimeJSON)
	if err := json.NewEncoder(w).Encode(&rsp); err != nil {
		log.Printf("%-7s %-7s cannot serialize fever JSON response: %s", logWarn, logName, err.Error())
	}

}
