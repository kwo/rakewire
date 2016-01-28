package httpd

import (
	"fmt"
	"net/http"
	"rakewire/model"
)

func statusHandler(w http.ResponseWriter, req *http.Request) {

	w.Header().Set(hContentType, mimeText)
	w.Write([]byte(fmt.Sprintf("Rakewire %s\n", model.Version)))
	w.Write([]byte(fmt.Sprintf("Build Time: %s\n", model.BuildTime)))
	w.Write([]byte(fmt.Sprintf("Build Hash: %s\n", model.BuildHash)))

}
