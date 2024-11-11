package decrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"os"
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

	// 복호화된 데이터를 저장할 버퍼
	plaintext := make([]byte, len(ciphertext))

	// CBC 블록 모드
	mode := cipher.NewCBCDecrypter(block, iv)
	// 블록 복호화 수행
	mode.CryptBlocks(plaintext, ciphertext)

	// todo 패딩이 있는지 확인 후 패딩 제거 실행 또는 패딩 제거 함수에서 패딩 여부 체크

	// PKCS7 패딩 제거
	unpaddedData, err := PKCS7UnPadding(plaintext)
	if err != nil {
		// 패딩이 없는 경우.
		if err.Error() == "잘못된 패딩" {
			return plaintext, nil
		} else { // 데이터가 비어있습니다.
			return nil, err
		}
	}
	// 패딩이 있는 경우.
	return unpaddedData, nil
}

func DecryptFile(filename string, key []byte) error {
	// 파일 읽기
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// 파일 복호화
	decryptedData, err := DecryptAES(data, key)
	if err != nil {
		return err
	}

	// 복호화된 데이터로 파일 덮어쓰기
	err = os.WriteFile(filename, decryptedData, 0644)
	if err != nil {
		return err
	}

	return nil
}

// PKCS7UnPadding 는 PKCS#7 패딩 제거 함수입니다.
func PKCS7UnPadding(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("데이터가 비어있습니다")
	}

	// 마지막 바이트가 패딩의 길이를 나타냄
	padding := int(data[len(data)-1])
	// 패딩 유효성 검사
	if padding > len(data) {
		return nil, errors.New("잘못된 패딩")
	}

	// 모든 패딩 바이트가 패딩 길이와 같은지 확인
	for i := len(data) - padding; i < len(data); i++ {
		if int(data[i]) != padding {
			return nil, errors.New("잘못된 패딩")
		}
	}

	// 패딩을 제외한 데이터 반환
	return data[:len(data)-padding], nil
}
