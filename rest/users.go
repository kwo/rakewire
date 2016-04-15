package rest

import (
	"encoding/json"
	"net/http"
)

func (z *API) usersGet(w http.ResponseWriter, req *http.Request) {

	rsp := &User{
		ID:       0,
		Username: "username",
	}

	w.Header().Set(hContentType, mimeJSON)
	if err := json.NewEncoder(w).Encode(&rsp); err != nil {
		log.Debugf("cannot serialize fever JSON response: %s", err.Error())
	}

}
