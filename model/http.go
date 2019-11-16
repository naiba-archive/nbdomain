package model

// Response ..
type Response struct {
	Code    uint64        `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}
