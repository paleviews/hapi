package runtime

import (
	"fmt"
)

func NewAPIError(code ResponseCode, desc string, src error) APIError {
	return APIError{
		Code:        code,
		Message:     desc,
		SourceError: src,
	}
}

type APIError struct {
	Code        ResponseCode `json:"code"`
	Message     string       `json:"message"`
	SourceError error        `json:"source_error,omitempty"`
}

func (e APIError) Error() string {
	s := fmt.Sprintf("hapi handler error code %d", e.Code)
	if e.Message != "" {
		s += " with message " + e.Message
	}
	if e.SourceError != nil {
		s += " and error " + e.SourceError.Error()
	}
	return s
}
