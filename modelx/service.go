package model

// Service standardizes the service interface.
type Service interface {
	Start() error
	Stop()
	IsRunning() bool
}
