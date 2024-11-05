package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// EncryptAES 는 data를 AES와 CBC 블록 모드로 암호화
func EncryptAES(data, key []byte) ([]byte, error) {
	if len(string(data))%aes.BlockSize != 0 {
		panic("plaintext is not a multiple of the block size")
	}

	// 키에 대한 블록 생성
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 암호문을 저장할 공간
	ciphertext := make([]byte, aes.BlockSize+len(data))
	// IV 초기화 벡터를 저장할 공간
	iv := ciphertext[:aes.BlockSize]
	// iv := []byte("abcdefghijklmnop")

	//iv에 랜덤값 설정
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// CBC 블록 모드
	mode := cipher.NewCBCEncrypter(block, iv)

	// 블록 암호화 수행
	mode.CryptBlocks(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil

}
