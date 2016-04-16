package httpd

import (
	"fmt"
	"net/http"
	"time"
)

func (z *Service) statusHandler(w http.ResponseWriter, req *http.Request) {

	w.Header().Set(hContentType, mimeText)
	w.Write([]byte(fmt.Sprintf("Rakewire %s\n", z.version)))
	w.Write([]byte(fmt.Sprintf("Uptime: %s\n", time.Now().Truncate(time.Second).Sub(z.appstart).String())))

}
