package httpd

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// http://feedafever.com/api

type feverResponse struct {
	XMLName       xml.Name  `json:"-" xml:"response"`
	Version       int       `json:"api_version" xml:"api_version"`
	Authorized    int       `json:"auth" xml:"auth"`
	LastRefreshed feverTime `json:"last_refreshed_on_time,omitempty" xml:"last_refreshed_on_time,omitempty"`
}

func (z *Service) feverRouter(prefix string) *mux.Router {

	router := mux.NewRouter()

	router.Queries("api", "").Methods(mPost).HandlerFunc(feverMux)
	router.Queries("api", "").HandlerFunc(notSupported)
	router.Path(prefix).HandlerFunc(notFound)

	return router

}

func feverMux(w http.ResponseWriter, req *http.Request) {

	if err := req.ParseForm(); err != nil {
		http.Error(w, "cannot parse request\n", 400)
		return
	}

	useXML := req.URL.Query().Get("api") == "xml"

	rsp := &feverResponse{
		Version: 3,
	}

	if apiKey := req.PostFormValue("api_key"); apiKey == "" {
		rsp.Authorized = 0
	} else {
		rsp.Authorized = 1
		// TODO: authenticate, api_key value is md5(email:password)
	}

	if rsp.Authorized == 1 {

		for k := range req.URL.Query() {

			switch k {
			case "api":
				rsp.LastRefreshed = feverTime{time.Now()} // TODO: get last refreshed
			case "feeds":
				// add to response
			}

		} // loop

	}

	if useXML {
		w.Header().Set(hContentType, "text/xml; charset=utf-8")
		w.Write([]byte(strings.ToLower(strings.TrimSuffix(xml.Header, "\n")))) // lowercase without trailing newline
		if err := xml.NewEncoder(w).Encode(&rsp); err != nil {
			log.Printf("%-7s %-7s cannot serialize fever XML response: %s", logWarn, logName, err.Error())
		}
	} else {
		w.Header().Set(hContentType, "text/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(&rsp); err != nil {
			log.Printf("%-7s %-7s cannot serialize fever JSON response: %s", logWarn, logName, err.Error())
		}
	}

}

type feverTime struct {
	time.Time
}

func (z feverTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", strconv.FormatInt(z.Unix(), 10))), nil
}

func (z feverTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	e.EncodeToken(start)
	e.EncodeToken(xml.CharData([]byte(strconv.FormatInt(z.Unix(), 10))))
	e.EncodeToken(xml.EndElement{Name: start.Name})
	return nil
}
