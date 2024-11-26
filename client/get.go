package client

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"math/big"
	"net/http"
)

// PublicKeyJSON 는 JSON에서 공개키로 변환하기 위한 구조체
type PublicKeyJSON struct {
	N string `json:"n"` // base64로 인코딩된 modulus
	E int    `json:"e"` // public exponent
}

func NewPublicKeyJSON() *PublicKeyJSON {
	return new(PublicKeyJSON)
}

// GetPublicKey 는 서버로 부터 공개키를 요청 합니다.
func (p *PublicKeyJSON) GetPublicKey(uri string) (*rsa.PublicKey, error) {

	// GET 요청 보내기
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	// body 닫기
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println("Error closing response body:", err)
		}
	}(resp.Body)

	// 응답 읽기
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// JSON 파싱
	var keyJSON PublicKeyJSON
	err = json.Unmarshal(body, &keyJSON)
	if err != nil {
		return nil, err
	}

	// RSA PublicKey로 변환
	publicKey, err := ConvertJSONToPublicKey(keyJSON)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("Received Public Key:\n")
	//fmt.Printf("N: %v\n", publicKey.N)
	//fmt.Printf("E: %v\n", publicKey.E)

	return publicKey, nil
}

// ConvertJSONToPublicKey 는 JSON을 rsa.PublicKey로 변환하는 함수
func ConvertJSONToPublicKey(keyJSON PublicKeyJSON) (*rsa.PublicKey, error) {
	// base64 디코딩
	nBytes, err := base64.StdEncoding.DecodeString(keyJSON.N)
	if err != nil {
		return nil, err
	}

	// rsa.PublicKey 구조체 생성
	pubKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: keyJSON.E,
	}

	return pubKey, nil
}
