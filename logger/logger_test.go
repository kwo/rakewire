package logger

import (
	"testing"
)

func TestInfo(t *testing.T) {
	log := New("tester")
	log.Infof("info message")
}

func TestDebug(t *testing.T) {
	log := New("tester")
	DebugMode = true
	log.Infof("msg 1 %d %d %d", 1, 2, 3)
	log.Debugf("msg 2 %d %d %d", 1, 2, 3)
	DebugMode = false
}

func TestLogger3(t *testing.T) {
	log := New("tester")
	log.Infof("info message %d %d %d", 1, 2, 3)
}

func TestLogger4(t *testing.T) {
	DateFormat = ""
	log := New("tester")
	log.Infof("info message")
}
