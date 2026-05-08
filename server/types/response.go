package types

type Response struct {
	Code int `json:"code"` // HTTP status code
	Data any `json:"data"`
}
