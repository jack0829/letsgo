package callback

import "time"

type status struct {
	Tries   int       `json:"tries,omitempty"`
	Error   string    `json:"error,omitempty"`
	ReqTime time.Time `json:"req_time"`
}
