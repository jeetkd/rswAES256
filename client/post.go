package client

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

// SendCipherWithPOST 는 서버로 암호문을 전송.
func SendCipherWithPOST(cipher []byte) error {
	// HTTP 요청 생성
	resp, err := http.Post(
		"http://localhost:8080/api/cipher",
		"application/octet-stream", //바이트 형식 전송.
		bytes.NewBuffer(cipher),
	)

	defer resp.Body.Close()

	if err != nil {
		log.Println("요청 전송 중 오류:", err)
		return fmt.Errorf("요청 전송 중 오류: %v", err)
	}
	return nil
}
