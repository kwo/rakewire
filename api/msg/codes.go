package msg

// Status codes
const (
	StatusOK       = 0
	StatusErr      = 1
	StatusNotFound = 2
)

// StatusText retrieves the status text for the given status code
func StatusText(code int) string {
	switch code {
	case StatusOK:
		return "OK"
	case StatusErr:
		return "Error"
	case StatusNotFound:
		return "Not Found"
	}
	return ""
}
