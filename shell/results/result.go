package results

import (
	"time"
)

// RestResponse -- Rest response
type RestResponse interface {
	GetContentType() string
	GetDuration() time.Duration
	GetPayloadDocument() ResultDocument
	GetHeaderDocument() ResultDocument
	GetCookieDocument() ResultDocument
}

type restResponse struct {
	docMap      map[string]ResultDocument
	duration    time.Duration
	contentType string
}
