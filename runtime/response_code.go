package runtime

type ResponseCode int32

// Default response codes are used if no proto enum is annotated as type of response code.
const (
	DefaultResponseCodeOK ResponseCode = iota
	DefaultResponseCodeInvalidInput
	DefaultResponseCodeUnauthenticated
	DefaultResponseCodeServerError
)

func GetDescFromResponseCode(code ResponseCode) string {
	switch code {
	case DefaultResponseCodeOK:
		return "ok"
	case DefaultResponseCodeInvalidInput:
		return "invalid_input"
	case DefaultResponseCodeUnauthenticated:
		return "unauthenticated"
	case DefaultResponseCodeServerError:
		return "server_error"
	default:
		return ""
	}
}

func APIErrorFromResponseCode(code ResponseCode, src error) APIError {
	return NewAPIError(code, GetDescFromResponseCode(code), src)
}
