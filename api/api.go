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
	handlers  map[string]map[string]Handler // handlers mapped by path then method
	version   string
	buildTime int64
	buildHash string
	appStart  int64
}

// Handler handles API requests
type Handler func(context.Context, http.ResponseWriter, *http.Request)

// New creates a new REST API instance
func New(database model.Database, mountPath, versionString string, appStart int64) *API {

	version, buildTime, buildHash := parseVersionString(versionString)

	z := &API{
		db:        database,
		mountPath: mountPath,
		handlers:  make(map[string]map[string]Handler),
		version:   version,
		buildTime: buildTime,
		buildHash: buildHash,
		appStart:  appStart,
	}

	// register handlers
	// TODO: handle more errRequest errors: auth

	z.handlers["groups/list"] = make(map[string]Handler)
	z.handlers["groups/list"][http.MethodPost] = func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		req := &msg.GroupListRequest{}
		if errRequest := readRequest(ctx, r, req); errRequest == nil {
			if rsp, errResponse := z.GroupList(ctx, req); errResponse == nil {
				sendResponse(ctx, w, rsp)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		} else if errRequest == ErrEmptyRequest {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}

	z.handlers["token"] = make(map[string]Handler)
	z.handlers["token"][http.MethodPost] = func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		req := &msg.TokenRequest{}
		if errRequest := readRequest(ctx, r, req); errRequest == nil {
			if rsp, errResponse := z.GetToken(ctx, req); errResponse == nil {
				sendResponse(ctx, w, rsp)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		} else if errRequest == ErrEmptyRequest {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}

	z.handlers["status"] = make(map[string]Handler)
	z.handlers["status"][http.MethodPost] = func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		req := &msg.StatusRequest{}
		if errRequest := readRequest(ctx, r, req); errRequest == nil {
			if rsp, errResponse := z.GetStatus(ctx, req); errResponse == nil {
				sendResponse(ctx, w, rsp)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		} else if errRequest == ErrEmptyRequest {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}

	z.handlers["subscriptions/add"] = make(map[string]Handler)
	z.handlers["subscriptions/add"][http.MethodPost] = func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		req := &msg.SubscriptionAddUpdateRequest{}
		if errRequest := readRequest(ctx, r, req); errRequest == nil {
			if rsp, errResponse := z.SubscriptionAddUpdate(ctx, req); errResponse == nil {
				sendResponse(ctx, w, rsp)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		} else if errRequest == ErrEmptyRequest {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}

	z.handlers["subscriptions/list"] = make(map[string]Handler)
	z.handlers["subscriptions/list"][http.MethodPost] = func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		req := &msg.SubscriptionListRequest{}
		if errRequest := readRequest(ctx, r, req); errRequest == nil {
			if rsp, errResponse := z.SubscriptionList(ctx, req); errResponse == nil {
				sendResponse(ctx, w, rsp)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		} else if errRequest == ErrEmptyRequest {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}

	z.handlers["subscriptions/unsubscribe"] = make(map[string]Handler)
	z.handlers["subscriptions/unsubscribe"][http.MethodPost] = func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		req := &msg.UnsubscribeRequest{}
		if errRequest := readRequest(ctx, r, req); errRequest == nil {
			if rsp, errResponse := z.SubscriptionUnsubscribe(ctx, req); errResponse == nil {
				sendResponse(ctx, w, rsp)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		} else if errRequest == ErrEmptyRequest {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}

	z.handlers["subscriptions.opml"] = make(map[string]Handler)
	z.handlers["subscriptions.opml"][http.MethodGet] = z.opmlExport
	z.handlers["subscriptions.opml"][http.MethodPut] = z.opmlImport

	return z

}

// ServeHTTPC context-aware http handler
func (z *API) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, z.mountPath)
	if handlers := z.handlers[path]; handlers != nil {
		if handler := handlers[r.Method]; handler != nil {
			handler(ctx, w, r)
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
		return ErrEmptyRequest
	}
	if errJSON := json.Unmarshal(data, req); errJSON != nil {
		return errJSON
	}
	return nil
}

func sendResponse(ctx context.Context, w http.ResponseWriter, rsp interface{}) {
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
