package client

import (
	"crypto/rsa"
	"rswAES256/config"
)

type Client struct {
	config    *config.Config
	PublicKey *rsa.PublicKey
	Key       string
}

// NewClient 는 서버로부터 키를 받아오기 위한 클라이언트 생성
func NewClient(config *config.Config) *Client {
	client := new(Client)

	client.config = config
	return client
}
