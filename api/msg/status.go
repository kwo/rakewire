package msg

// StatusRequest defines the status request
type StatusRequest struct{}

// StatusResponse defines the status response
type StatusResponse struct {
	Version   string `json:"version,omitempty"`
	BuildTime int64  `json:"buildTime,omitempty"`
	BuildHash string `json:"buildHash,omitempty"`
	AppStart  int64  `json:"appStart,omitempty"`
}
