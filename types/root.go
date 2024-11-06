package types

// KeyRequest 는 서버로부터 받아올 KEY
type KeyRequest struct {
	Key string `json:"key"`
}
