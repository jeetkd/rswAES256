package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// JSON으로 변환을 위한 구조체
type PublicKeyJSON struct {
	N string `json:"n"` // base64로 인코딩된 modulus
	E int    `json:"e"` // public exponent
}

// RSA 공개키를 JSON 형식으로 변환하는 함수
func ConvertPublicKeyToJSON(key *rsa.PublicKey) PublicKeyJSON {
	return PublicKeyJSON{
		N: base64.StdEncoding.EncodeToString(key.N.Bytes()),
		E: key.E,
	}
}

// HTTP 핸들러 함수
func GetPublicKeyHandler(chPKey chan *rsa.PrivateKey) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// RSA 키 쌍 생성
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			http.Error(w, "Failed to generate key", http.StatusInternalServerError)
			return
		}
		fmt.Println("N : ", privateKey.PublicKey.N)
		fmt.Println("E : ", privateKey.PublicKey.E)

		// 공개키를 JSON 형식으로 변환
		publicKeyJSON := ConvertPublicKeyToJSON(&privateKey.PublicKey)
		fmt.Println("N2 :", publicKeyJSON.N)
		fmt.Println("E2 :", publicKeyJSON.E)
		// JSON 응답 설정
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(publicKeyJSON)
		chPKey <- privateKey
	}
}

// POST 요청을 처리하는 핸들러 함수
func PostCipherHandler(cipher chan []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// POST 메서드 확인
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 요청 본문 읽기
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// 받은 바이트 데이터 처리
		//fmt.Printf("Received byte data: %v\n", body)
		//fmt.Printf("Received string: %s\n", string(body))

		cipher <- body
	}

}

func main() {
	chPKey := make(chan *rsa.PrivateKey)
	cipher := make(chan []byte)
	// GET : 클라이언트가 요청시 공개키 전송.
	http.HandleFunc("/api/publickey", GetPublicKeyHandler(chPKey))
	// POST : 클라이언트가 데이터 전송시 암호문을 가져옴.
	http.HandleFunc("/api/cipher", PostCipherHandler(cipher))

	// 서버 시작
	go http.ListenAndServe(":8080", nil)

	recvPKey := <-chPKey
	recvCipher := <-cipher

	plaintext, err := rsa.DecryptPKCS1v15( // 암호화된 데이터를 개인 키로 복호화
		rand.Reader,
		recvPKey, // 개인키
		recvCipher,
	)
	if err != nil {
		panic(err)
	}
	// 개인키로 암호화된 데이터 복호화후 를 출력
	fmt.Printf("plantext : %x", string(plaintext))
}
