package api

import (
	"encoding/json"
	"github.com/kwo/rakewire/api/msg"
	"github.com/kwo/rakewire/logger"
	"github.com/kwo/rakewire/model"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	hContentType = "Content-Type"
	mimeJSON     = "application/json"
)

var (
	log = logger.New("api")
)

// API top level struct
type API struct {
	db        model.Database
	mountPath string
	version   string
	buildTime int64
	buildHash string
	appStart  int64
}

// New creates a new REST API instance
func New(database model.Database, mountPath, versionString string, appStart int64) *API {

	version, buildTime, buildHash := parseVersionString(versionString)

	return &API{
		db:        database,
		mountPath: mountPath,
		version:   version,
		buildTime: buildTime,
		buildHash: buildHash,
		appStart:  appStart,
	}

}

// ServeHTTPC context-aware http handler
func (z *API) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, z.mountPath)
	if path == "token" {
		if r.Method == http.MethodPost {
			req := &msg.TokenRequest{}
			if errRequest := readRequest(ctx, r, req); errRequest == nil {
				rsp, errToken := z.GetToken(ctx, req)
				sendResponse(ctx, w, rsp, errToken)
			} else {
				sendResponse(ctx, w, nil, errRequest)
			}
		} else {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	} else if path == "status" {
		if r.Method == http.MethodPost {
			req := &msg.StatusRequest{}
			if errRequest := readRequest(ctx, r, req); errRequest == nil {
				rsp, errStatus := z.GetStatus(ctx, &msg.StatusRequest{})
				sendResponse(ctx, w, rsp, errStatus)
			} else {
				sendResponse(ctx, w, nil, errRequest)
			}
		} else {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	} else if path == "subscriptions.opml" {
		if r.Method == http.MethodGet {
			z.opmlExport(ctx, w, r)
		} else if r.Method == http.MethodPut {
			z.opmlImport(ctx, w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func readRequest(ctx context.Context, r *http.Request, req interface{}) error {
	defer r.Body.Close()
	data, errRead := ioutil.ReadAll(r.Body)
	if errRead != nil {
		return errRead
	}
	if len(data) == 0 {
		return msg.ErrEmptyRequest
	}
	if errJSON := json.Unmarshal(data, req); errJSON != nil {
		return errJSON
	}
	return nil
}

func sendResponse(ctx context.Context, w http.ResponseWriter, rsp interface{}, errRequest error) {
	if errRequest == nil {
		data, errJSON := json.Marshal(rsp)
		if errJSON == nil {
			w.Header().Set(hContentType, mimeJSON)
			data = append(data, '\n')
			if _, err := w.Write(data); err != nil {
				log.Debugf("Error sending response: %s", err.Error())
			}
		} else {
			log.Debugf("Cannot serialize response: %s", errJSON.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	} else if errRequest == msg.ErrEmptyRequest {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	} else {
		// TODO: senteniel errors for triage
		// not authorized
		// bad request
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func parseVersionString(versionString string) (string, int64, string) {

	// parse version string
	fields := strings.Fields(versionString)
	if len(fields) == 3 {

		version := fields[0]
		buildHash := fields[2]

		buildTime, err := time.Parse(time.RFC3339, fields[1])
		if err != nil {
			log.Debugf("Cannot parse build time: %s", err.Error())
		}

		return version, buildTime.Unix(), buildHash

	}

	return "", 0, ""

}
