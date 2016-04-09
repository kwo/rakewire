package httpd

import (
	"fmt"
	"net/http"
	"rakewire/model"
	"time"
)

func statusHandler(w http.ResponseWriter, req *http.Request) {

	uptimeString := time.Now().Truncate(time.Second).Sub(model.AppStart.Truncate(time.Second)).String()

	w.Header().Set(hContentType, mimeText)
	w.Write([]byte(fmt.Sprintf("Rakewire %s\n", model.Version)))
	w.Write([]byte(fmt.Sprintf("Commit Time: %s\n", model.CommitTime)))
	w.Write([]byte(fmt.Sprintf("Commit Hash: %s\n", model.CommitHash)))
	w.Write([]byte(fmt.Sprintf("Uptime: %s\n", uptimeString)))

}
