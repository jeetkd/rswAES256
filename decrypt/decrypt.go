package decrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

func DecryptAES(ciphertext, key []byte) ([]byte, error) {
	// 암호문의 길이가 최소한 IV 크기보다는 커야 함
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("암호문이 너무 짧습니다")
	}

	// 암호화 블록 생성
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// IV 추출 (암호문의 첫 BlockSize 바이트)
	iv := ciphertext[:aes.BlockSize]

	// 실제 암호문
	ciphertext = ciphertext[aes.BlockSize:]

	// CBC 모드는 항상 전체 블록 단위로 작동합니다
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("암호문이 블록 크기의 배수가 아닙니다")
	}

	// 복호화된 데이터를 저장할 버퍼
	plaintext := make([]byte, len(ciphertext))

	// CBC 블록 모드
	mode := cipher.NewCBCDecrypter(block, iv)

	// 블록 복호화 수행
	mode.CryptBlocks(plaintext, ciphertext)

	return plaintext, nil
}
