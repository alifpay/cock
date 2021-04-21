package models

type Response struct {
	ExternalRef string      `json:"external_ref"`
	Code        int         `json:"code"`
	Message     string      `json:"message"`
	Payload     interface{} `json:"payload"`
	ID          uint64      `json:"id"`
}
